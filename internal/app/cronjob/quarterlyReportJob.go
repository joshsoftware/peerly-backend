package cronjob

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	apprSvc "github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	"github.com/joshsoftware/peerly-backend/internal/app/googlesheets"
	reportSvc "github.com/joshsoftware/peerly-backend/internal/app/reportAppreciations"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

const QUARTERLY_REPORT_JOB = "QUARTERLY_REPORT_JOB"

// 8:25 AM IST = 2:55 AM UTC → run on 1st day of Jan, Apr, Jul, Oct
const QUARTERLY_CRON_EXPRESSION = "55 2 1 1,4,7,10 *"

var quarterMonthNames = map[int]string{
	1: "Mar-May",
	2: "Jun-Aug",
	3: "Sep-Nov",
	4: "Dec-Feb",
}

type QuarterlyReportJob struct {
	CronJob
	appreciationService       apprSvc.Service
	reportAppreciationService reportSvc.Service
	sheetService              *googlesheets.Service
	spreadsheetID             string
}

func NewQuarterlyReportJob(
	appreciationSvc apprSvc.Service,
	reportAppreciationSvc reportSvc.Service,
	sheetSvc *googlesheets.Service,
	spreadsheetID string,
	scheduler gocron.Scheduler,
) Job {
	return &QuarterlyReportJob{
		appreciationService:       appreciationSvc,
		reportAppreciationService: reportAppreciationSvc,
		sheetService:              sheetSvc,
		spreadsheetID:             spreadsheetID,
		CronJob: CronJob{
			name:      QUARTERLY_REPORT_JOB,
			scheduler: scheduler,
		},
	}
}

func (cron *QuarterlyReportJob) Schedule() error {
	var err error
	cron.job, err = cron.scheduler.NewJob(
		gocron.CronJob(QUARTERLY_CRON_EXPRESSION, false),
		gocron.NewTask(cron.Execute, cron.Task),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	cron.scheduler.Start()
	if err != nil {
		logger.Warn(context.TODO(), fmt.Sprintf("error occurred while scheduling %s, message %+v", cron.name, err.Error()))
		return err
	}
	logger.Info(context.TODO(), fmt.Sprintf("Quarterly report job scheduled (cron: %s)", QUARTERLY_CRON_EXPRESSION))
	return nil
}

func (cron *QuarterlyReportJob) Task(ctx context.Context) {
	logger.Info(ctx, "in quarterly report job task")

	var err error
	for i := 0; i < 3; i++ {
		logger.Infof(ctx, "quarterly report job attempt: %d", i+1)
		err = cron.exportQuarterlyReport(ctx)
		if err == nil {
			logger.Info(ctx, "quarterly report exported successfully to Google Sheet")
			return
		}
		logger.Errorf(ctx, "quarterly report job attempt %d failed: %v", i+1, err)
	}
	logger.Errorf(ctx, "quarterly report job failed after 3 attempts: %v", err)
}

func (cron *QuarterlyReportJob) exportQuarterlyReport(ctx context.Context) error {

	quarter, year := getPreviousQuarterAndYear()
	logger.Infof(ctx, "Exporting quarterly report for Q%d(%d) to Google Sheet", quarter, year)

	quarterStart, quarterEnd := user.GetQuarterRangeUnixTime(quarter, year)

	// Inject a dummy user ID into the context to satisfy the API's token parsing requirement
	ctx = context.WithValue(ctx, constants.UserId, int64(0))

	// Fetch all appreciations
	filter := dto.AppreciationFilter{
		Self:  false,
		Limit: constants.DefaultPageSize,
		Page:  1,
	}
	appreciationResp, err := cron.appreciationService.ListAppreciations(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to list appreciations: %w", err)
	}

	var quarterAppreciations []dto.AppreciationResponse
	for _, appr := range appreciationResp.Appreciations {
		if appr.CreatedAt >= quarterStart && appr.CreatedAt < quarterEnd {
			quarterAppreciations = append(quarterAppreciations, appr)
		}
	}
	logger.Infof(ctx, "Found %d appreciations for Q%d(%d)", len(quarterAppreciations), quarter, year)

	// Fetch reported appreciations and build a lookup map by appreciation_id
	reportedMap := make(map[int64]dto.ReportedAppreciation)
	reportedResp, err := cron.reportAppreciationService.ListReportedAppreciations(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to list reported appreciations (continuing without report data): %v", err)
	} else {
		for _, reported := range reportedResp.Appreciations {
			reportedMap[reported.Appreciation_id] = reported
		}
	}

	tabName := fmt.Sprintf("Q%d(%d) %s", quarter, year, quarterMonthNames[quarter])
	rows := buildSheetRows(quarterAppreciations, reportedMap)

	err = cron.sheetService.CreateTab(cron.spreadsheetID, tabName)
	if err != nil {
		return fmt.Errorf("failed to create tab: %w", err)
	}

	err = cron.sheetService.AppendRows(cron.spreadsheetID, tabName, rows)
	if err != nil {
		return fmt.Errorf("failed to append rows: %w", err)
	}

	logger.Infof(ctx, "Successfully exported %d appreciations to tab '%s'", len(quarterAppreciations), tabName)
	return nil
}

func getPreviousQuarterAndYear() (quarter int, year int) {
	ist, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(ist)
	month := now.Month()
	y := now.Year()

	switch {
	case month >= time.July && month < time.October:
		return 1, y // Q1 (Mar-May) of current year
	case month >= time.October && month <= time.December:
		return 2, y // Q2 (Jun-Aug) of current year
	case month >= time.January && month < time.April:
		return 3, y - 1 // Q3 (Sep-Nov) of previous year
	case month >= time.April && month < time.July:
		return 4, y - 1 // Q4 (Dec-Feb) of previous year
	}
	return 1, y
}

// buildSheetRows constructs the header + data rows matching the 21-column sheet format.
func buildSheetRows(appreciations []dto.AppreciationResponse, reportedMap map[int64]dto.ReportedAppreciation) [][]interface{} {
	headers := []interface{}{
		"Core value", "Time Stamp", "Sender Employee ID", "Sender Full Name",
		"Sender designation", "Receiver Employee ID", "Receiver Full Name",
		"Receiver designation", "Total rewards", "Total reward points",
		"Appreciated Date", "Quarter", "Year", "Reporter Emp ID",
		"Reporting Comment", "Reporter Full Name", "Reported Date",
		"Moderator comment", "Moderator Full Name", "Moderated Date",
		"Appreciation Status",
	}

	rows := [][]interface{}{headers}

	for _, appr := range appreciations {
		createdTime := time.UnixMilli(appr.CreatedAt)
		appreciatedDate := createdTime.Format("02/01/2006")
		quarterName := user.GetQuarterName(createdTime)

		reporterEmpID := ""
		reportingComment := ""
		reporterFullName := ""
		reportedDate := ""
		moderatorComment := ""
		moderatorFullName := ""
		moderatedDate := ""
		appreciationStatus := ""

		if reported, ok := reportedMap[appr.ID]; ok {
			reporterEmpID = reported.ReporterEmployeeID
			reportingComment = reported.ReportingComment
			reporterFullName = reported.ReportedByFirstName + " " + reported.ReportedByLastName
			if reported.ReportedAt > 0 {
				reportedDate = time.UnixMilli(reported.ReportedAt).Format("02/01/2006")
			}
			moderatorComment = reported.ModeratorComment
			if reported.ModeratedByFirstName != "" || reported.ModeratedByLastName != "" {
				moderatorFullName = reported.ModeratedByFirstName + " " + reported.ModeratedByLastName
			}
			if reported.ModeratedAt > 0 {
				moderatedDate = time.UnixMilli(reported.ModeratedAt).Format("02/01/2006")
			}
			appreciationStatus = reported.Status
		}

		row := []interface{}{
			appr.CoreValueName,    // Core value
			appreciatedDate,       // Time Stamp
			appr.SenderEmployeeID, // Sender Employee ID
			appr.SenderFirstName + " " + appr.SenderLastName,     // Sender Full Name
			appr.SenderDesignation,                               // Sender designation
			appr.ReceiverEmployeeID,                              // Receiver Employee ID
			appr.ReceiverFirstName + " " + appr.ReceiverLastName, // Receiver Full Name
			appr.ReceiverDesignation,                             // Receiver designation
			appr.TotalRewards,                                    // Total rewards
			appr.TotalRewardPoints,                               // Total reward points
			appreciatedDate,                                      // Appreciated Date
			quarterName,                                          // Quarter
			createdTime.Year(),                                   // Year
			reporterEmpID,                                        // Reporter Emp ID
			reportingComment,                                     // Reporting Comment
			reporterFullName,                                     // Reporter Full Name
			reportedDate,                                         // Reported Date
			moderatorComment,                                     // Moderator comment
			moderatorFullName,                                    // Moderator Full Name
			moderatedDate,                                        // Moderated Date
			appreciationStatus,                                   // Appreciation Status
		}

		rows = append(rows, row)
	}

	return rows
}
