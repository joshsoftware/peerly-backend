package reportappreciations

import (
	"context"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func TestReportAppreciation(t *testing.T) {
	reportAppreciationRepo := mocks.NewReportAppreciationStorer(t)
	userRepo := mocks.NewUserStorer(t)
	service := NewService(reportAppreciationRepo, userRepo)

	tests := []struct {
		name            string
		userId          int64
		reqData         dto.ReportAppreciationReq
		setup           func(reportAppreciationMock *mocks.ReportAppreciationStorer)
		isErrorExpected bool
	}{
		{
			name:   "Success for report appreciation",
			userId: 1334,
			reqData: dto.ReportAppreciationReq{
				ReportingComment: "reporting comment",
				AppreciationId:   4,
			},
			setup: func(reportAppreciationMock *mocks.ReportAppreciationStorer) {
				reportAppreciationMock.On("CheckAppreciation", mock.Anything, mock.Anything).Return(true, nil).Once()
				reportAppreciationMock.On("CheckDuplicateReport", mock.Anything, mock.Anything).Return(false, nil).Once()
				reportAppreciationMock.On("GetSenderAndReceiver", mock.Anything, mock.Anything).Return(dto.GetSenderAndReceiverResp{
					Sender:   1004,
					Receiver: 1100,
				}, nil).Once()
				reportAppreciationMock.On("ReportAppreciation", mock.Anything, mock.Anything).Return(dto.ReportAppricaitionResp{}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:   "Error in report appreciation",
			userId: 1334,
			reqData: dto.ReportAppreciationReq{
				ReportingComment: "reporting comment",
				AppreciationId:   4,
			},
			setup: func(reportAppreciationMock *mocks.ReportAppreciationStorer) {
				reportAppreciationMock.On("CheckAppreciation", mock.Anything, mock.Anything).Return(true, nil).Once()
				reportAppreciationMock.On("CheckDuplicateReport", mock.Anything, mock.Anything).Return(false, nil).Once()
				reportAppreciationMock.On("GetSenderAndReceiver", mock.Anything, mock.Anything).Return(dto.GetSenderAndReceiverResp{
					Sender:   1004,
					Receiver: 1100,
				}, nil).Once()
				reportAppreciationMock.On("ReportAppreciation", mock.Anything, mock.Anything).Return(dto.ReportAppricaitionResp{}, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
		{
			name:   "Sender cannot report its own appreciation",
			userId: 1334,
			reqData: dto.ReportAppreciationReq{
				ReportingComment: "reporting comment",
				AppreciationId:   4,
			},
			setup: func(reportAppreciationMock *mocks.ReportAppreciationStorer) {
				reportAppreciationMock.On("CheckAppreciation", mock.Anything, mock.Anything).Return(true, nil).Once()
				reportAppreciationMock.On("CheckDuplicateReport", mock.Anything, mock.Anything).Return(false, nil).Once()
				reportAppreciationMock.On("GetSenderAndReceiver", mock.Anything, mock.Anything).Return(dto.GetSenderAndReceiverResp{
					Sender:   1334,
					Receiver: 1100,
				}, nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:   "Receiver cannot report its own appreciation",
			userId: 1334,
			reqData: dto.ReportAppreciationReq{
				ReportingComment: "reporting comment",
				AppreciationId:   4,
			},
			setup: func(reportAppreciationMock *mocks.ReportAppreciationStorer) {
				reportAppreciationMock.On("CheckAppreciation", mock.Anything, mock.Anything).Return(true, nil).Once()
				reportAppreciationMock.On("CheckDuplicateReport", mock.Anything, mock.Anything).Return(false, nil).Once()
				reportAppreciationMock.On("GetSenderAndReceiver", mock.Anything, mock.Anything).Return(dto.GetSenderAndReceiverResp{
					Sender:   1100,
					Receiver: 1334,
				}, nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:   "Duplicate report",
			userId: 1334,
			reqData: dto.ReportAppreciationReq{
				ReportingComment: "reporting comment",
				AppreciationId:   4,
			},
			setup: func(reportAppreciationMock *mocks.ReportAppreciationStorer) {
				reportAppreciationMock.On("CheckAppreciation", mock.Anything, mock.Anything).Return(true, nil).Once()
				reportAppreciationMock.On("CheckDuplicateReport", mock.Anything, mock.Anything).Return(true, nil).Once()
			},
			isErrorExpected: true,
		},
		{
			name:   "Appreciation does not exist",
			userId: 1334,
			reqData: dto.ReportAppreciationReq{
				ReportingComment: "reporting comment",
				AppreciationId:   4,
			},
			setup: func(reportAppreciationMock *mocks.ReportAppreciationStorer) {
				reportAppreciationMock.On("CheckAppreciation", mock.Anything, mock.Anything).Return(false, nil).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, constants.UserId, int64(test.userId))
			test.setup(reportAppreciationRepo)

			// test service
			_, err := service.ReportAppreciation(ctx, test.reqData)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestListReportedAppreciations(t *testing.T) {
	reportAppreciationRepo := mocks.NewReportAppreciationStorer(t)
	userRepo := mocks.NewUserStorer(t)
	service := NewService(reportAppreciationRepo, userRepo)

	tests := []struct {
		name            string
		ctx             context.Context
		setup           func(reportAppreciationMock *mocks.ReportAppreciationStorer)
		isErrorExpected bool
	}{
		{
			name: "Success for report appreciation",
			setup: func(reportAppreciationMock *mocks.ReportAppreciationStorer) {
				reportAppreciationMock.On("ListReportedAppreciations", mock.Anything).Return([]repository.ListReportedAppreciations{}, nil).Once()
			},
			isErrorExpected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, constants.UserId, int64(test.userId))
			test.setup(reportAppreciationRepo)

			// test service
			_, err := service.ReportAppreciation(ctx, test.reqData)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}
