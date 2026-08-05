package cronjob

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-co-op/gocron/v2"
	apprSvc "github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	"github.com/joshsoftware/peerly-backend/internal/app/email"
	"github.com/joshsoftware/peerly-backend/internal/app/googlesheets"
	reportSvc "github.com/joshsoftware/peerly-backend/internal/app/reportAppreciations"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"google.golang.org/api/sheets/v4"
)

const QUARTERLY_REPORT_JOB = "QUARTERLY_REPORT_JOB"

// 8:25 AM IST → run on 1st day of Jan, Apr, Jul, Oct
const QUARTERLY_CRON_EXPRESSION = "25 8 1 1,4,7,10 *"

const ALL_APPRECIATIONS_TAB_NAME = "All Appreciations"

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
		logger.CronWarn(context.TODO(), fmt.Sprintf("error occurred while scheduling %s, message %+v", cron.name, err.Error()))
		return err
	}
	logger.CronInfo(context.TODO(), fmt.Sprintf("Quarterly report job scheduled (cron: %s)", QUARTERLY_CRON_EXPRESSION))
	return nil
}

func (cron *QuarterlyReportJob) Task(ctx context.Context) {
	logger.CronInfo(ctx, "in quarterly report job task")

	var err error
	for i := 0; i < 3; i++ {
		logger.CronInfof(ctx, "quarterly report job attempt: %d", i+1)
		err = cron.ExportQuarterlyReport(ctx)
		if err == nil {
			logger.CronInfo(ctx, "quarterly report exported successfully to Google Sheet")
			return
		}
		logger.CronErrorf(ctx, "quarterly report job attempt %d failed: %v", i+1, err)
	}
	logger.CronErrorf(ctx, "quarterly report job failed after 3 attempts: %v", err)
}

func (cron *QuarterlyReportJob) ExportQuarterlyReport(ctx context.Context) error {

	quarter, year := getPreviousQuarterAndYear()
	logger.CronInfof(ctx, "Exporting quarterly report for Q%d(%d) to Google Sheet", quarter, year)

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
	logger.CronInfof(ctx, "Found %d appreciations for Q%d(%d)", len(quarterAppreciations), quarter, year)

	// Fetch reported appreciations and build a lookup map by appreciation_id
	reportedMap := make(map[int64]dto.ReportedAppreciation)
	reportedResp, err := cron.reportAppreciationService.ListReportedAppreciations(ctx, quarter, year)
	if err != nil {
		logger.CronErrorf(ctx, "failed to list reported appreciations (continuing without report data): %v", err)
	} else {
		for _, reported := range reportedResp.Appreciations {
			reportedMap[reported.Appreciation_id] = reported
		}
	}

	// Sort the Appreciations slice by CreatedAt in ascending order
	sort.SliceStable(quarterAppreciations, func(i, j int) bool {
		return quarterAppreciations[i].CreatedAt < quarterAppreciations[j].CreatedAt
	})

	// tabName := getFinancialYearTabName(year)
	tabName := ALL_APPRECIATIONS_TAB_NAME
	allRows := buildSheetRows(quarterAppreciations, reportedMap)

	tabExists, err := cron.sheetService.TabExists(cron.spreadsheetID, tabName)
	if err != nil {
		logger.CronErrorf(ctx, "failed to check if tab exists: %v", err)
	}

	if !tabExists {
		err = cron.sheetService.CreateTab(cron.spreadsheetID, tabName)
		if err != nil {
			return fmt.Errorf("failed to create tab: %w", err)
		}

		err = cron.sheetService.AppendRows(cron.spreadsheetID, tabName, allRows)
		if err != nil {
			return fmt.Errorf("failed to append rows: %w", err)
		}

		headerLength := 0
		if len(allRows) > 0 {
			headerLength = len(allRows[0])
		}

		formatParams := googlesheets.HeaderFormatParams{
			SpreadsheetID: cron.spreadsheetID,
			TabName:       tabName,
			HeaderLength:  headerLength,
			Color1: &sheets.Color{
				Red:   1.0,
				Green: 0.95,
				Blue:  0.8,
			},
			Color2: &sheets.Color{
				Red:   0.85,
				Green: 0.92,
				Blue:  0.83,
			},
			Interval1: 5,
			Interval2: 3,
		}

		// Apply header formatting
		err = cron.sheetService.FormatHeaderRow(formatParams)
		if err != nil {
			logger.CronErrorf(ctx, "failed to format header row: %v", err)
		}
	} else {
		var dataRows [][]interface{}
		if len(allRows) > 1 {
			dataRows = allRows[1:]
		}
		if len(dataRows) > 0 {
			err = cron.sheetService.AppendRows(cron.spreadsheetID, tabName, dataRows)
			if err != nil {
				return fmt.Errorf("failed to append data rows: %w", err)
			}
		} else {
			logger.CronInfof(ctx, "No data rows to append for tab '%s'", tabName)
		}
	}

	logger.CronInfof(ctx, "Successfully exported %d appreciations to tab '%s'", len(quarterAppreciations), tabName)

	sheetLink := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit", cron.spreadsheetID)
	cron.sendSuccessEmail(ctx, sheetLink)

	return nil
}

type EmailData struct {
	SheetLink string
}

func (cron *QuarterlyReportJob) sendSuccessEmail(ctx context.Context, sheetLink string) {
	subject := "Peerly Automated Report For Looker"
	mailReq := email.NewMail([]string{constants.DL["hr_team"]}, []string{}, []string{}, subject)

	data := EmailData{
		SheetLink: sheetLink,
	}

	err := mailReq.ParseTemplate("internal/app/email/templates/quarterlyReport.html", data)
	if err != nil {
		logger.CronErrorf(ctx, "Failed to parse email template: %v", err)
		return
	}

	err = mailReq.Send()
	if err != nil {
		logger.CronErrorf(ctx, "Failed to send quarterly report email: %v", err)
	} else {
		logger.CronInfo(ctx, "Successfully sent quarterly report email to HR team")
	}
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

// func getFinancialYearTabName(year int) string {
// 	nextYear := (year + 1) % 100
// 	return fmt.Sprintf("FY %d-%02d", year, nextYear)
// }

// buildSheetRows constructs the header + data rows matching the 21-column sheet format.
func buildSheetRows(appreciations []dto.AppreciationResponse, reportedMap map[int64]dto.ReportedAppreciation) [][]interface{} {
	headers := []interface{}{
		"Core value", "Time Stamp", "Sender Employee ID", "Sender Full Name",
		"Sender designation", "Receiver Employee ID", "Receiver Full Name",
		"Receiver designation", "Appreciation description", "Total rewards", "Total reward points",
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
			appr.Description,                                     // Appreciation description
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
