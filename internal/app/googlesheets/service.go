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
