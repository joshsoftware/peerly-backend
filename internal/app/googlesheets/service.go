package googlesheets

import (
	"context"
	"fmt"
	"os"

	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Service wraps the Google Sheets API client.
type Service struct {
	sheetsService *sheets.Service
}

// NewService creates a new Google Sheets service using the given service account credentials file.
func NewService(credentialsPath string) (*Service, error) {
	ctx := context.Background()
	creds, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read google credentials file: %w", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	return &Service{sheetsService: srv}, nil
}

// TabExists checks if a sheet tab exists in the spreadsheet.
func (s *Service) TabExists(spreadsheetID, tabName string) (bool, error) {
	spreadsheet, err := s.sheetsService.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return false, fmt.Errorf("failed to get spreadsheet details: %w", err)
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == tabName {
			return true, nil
		}
	}
	return false, nil
}

// CreateTab creates a new sheet tab in the spreadsheet. If the tab already exists, it silently succeeds.
func (s *Service) CreateTab(spreadsheetID, tabName string) error {
	spreadsheet, err := s.sheetsService.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return fmt.Errorf("failed to get spreadsheet details: %w", err)
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == tabName {
			logger.Info(context.Background(), "Google Sheets: tab already exists: ", tabName)
			return nil
		}
	}

	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: tabName,
					},
				},
			},
		},
	}

	_, err = s.sheetsService.Spreadsheets.BatchUpdate(spreadsheetID, req).Do()
	if err != nil {
		return fmt.Errorf("failed to create tab '%s': %w", tabName, err)
	}

	logger.Info(context.Background(), "Google Sheets: created tab: ", tabName)
	return nil
}

// AppendRows appends the given rows (including header) to the specified tab.
func (s *Service) AppendRows(spreadsheetID, tabName string, rows [][]interface{}) error {
	valueRange := &sheets.ValueRange{
		Values: rows,
	}

	_, err := s.sheetsService.Spreadsheets.Values.Append(
		spreadsheetID,
		tabName,
		valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		return fmt.Errorf("failed to append rows to tab '%s': %w", tabName, err)
	}

	logger.Infof(context.Background(), "Google Sheets: appended %d rows to tab '%s'", len(rows), tabName)
	return nil
}

// getSheetID finds the Sheet ID for a given tab name.
func (s *Service) getSheetID(spreadsheetID, tabName string) (int64, error) {
	spreadsheet, err := s.sheetsService.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return 0, err
	}
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == tabName {
			return sheet.Properties.SheetId, nil
		}
	}
	return 0, fmt.Errorf("tab '%s' not found", tabName)
}

// HeaderFormatParams holds configuration for alternating header colors.
type HeaderFormatParams struct {
	SpreadsheetID string
	TabName       string
	HeaderLength  int
	Color1        *sheets.Color
	Color2        *sheets.Color
	Interval1     int
	Interval2     int
}

// FormatHeaderRow applies alternating background colors to the first row (header) of the tab.
func (s *Service) FormatHeaderRow(params HeaderFormatParams) error {
	sheetID, err := s.getSheetID(params.SpreadsheetID, params.TabName)
	if err != nil {
		return err
	}

	var requests []*sheets.Request

	col := 0
	useColor1 := true

	for col < params.HeaderLength {
		var interval int
		var color *sheets.Color

		if useColor1 {
			interval = params.Interval1
			color = params.Color1
		} else {
			interval = params.Interval2
			color = params.Color2
		}

		endCol := col + interval
		if endCol > params.HeaderLength {
			endCol = params.HeaderLength
		}

		requests = append(requests, &sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartRowIndex:    0,
					EndRowIndex:      1, // Only the first row
					StartColumnIndex: int64(col),
					EndColumnIndex:   int64(endCol),
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: color,
					},
				},
				Fields: "userEnteredFormat(backgroundColor)",
			},
		})

		col = endCol
		useColor1 = !useColor1
	}

	if len(requests) == 0 {
		return nil
	}

	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	_, err = s.sheetsService.Spreadsheets.BatchUpdate(params.SpreadsheetID, req).Do()
	if err != nil {
		return fmt.Errorf("failed to format header for tab '%s': %w", params.TabName, err)
	}

	logger.Info(context.Background(), "Google Sheets: formatted header row for tab: ", params.TabName)
	return nil
}
