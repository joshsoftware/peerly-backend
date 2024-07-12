package repository

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type ReportAppreciationStorer interface {
	ReportAppreciation(ctx context.Context, reportReq dto.ReportAppreciationReq) (resp dto.ReportAppricaitionResp, err error)
	GetSenderAndReceiver(ctx context.Context, reqData dto.ReportAppreciationReq) (resp dto.GetSenderAndReceiverResp, err error)
	CheckDuplicateReport(ctx context.Context, reqData dto.ReportAppreciationReq) (isDupliate bool, err error)
	CheckAppreciation(ctx context.Context, reqData dto.ReportAppreciationReq) (doesExist bool, err error)
}

type Resolution struct {
	Id               int64  `json:"id" db:"id"`
	AppreciationId   int64  `json:"appreciation_id" db:"appreciation_id"`
	ReportingComment string `json:"reporting_comment" db:"reporting_comment"`
	ReportedBy       int64  `json:"reported_by" db:"reported_by"`
	ReportedAt       int64  `json:"reported_at" db:"reported_at"`
	ModeratorAction  int64  `json:"moderator_action" db:"moderator_action"`
	ModeratorComment string `json:"moderator_comment" db:"moderator_comment"`
	ModeratedBy      int64  `json:"moderated_by" db:"moderated_by"`
	ModeratedAt      int64  `json:"moderated_at" db:"moderated_at"`
}
