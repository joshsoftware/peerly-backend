package repository

import (
	"context"
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type ReportAppreciationStorer interface {
	ReportAppreciation(ctx context.Context, reportReq dto.ReportAppreciationReq) (resp dto.ReportAppricaitionResp, err error)
	GetSenderAndReceiver(ctx context.Context, reqData dto.ReportAppreciationReq) (resp dto.GetSenderAndReceiverResp, err error)
	CheckDuplicateReport(ctx context.Context, reqData dto.ReportAppreciationReq) (isDupliate bool, err error)
	CheckAppreciation(ctx context.Context, reqData dto.ReportAppreciationReq) (doesExist bool, err error)
	ListReportedAppreciations(ctx context.Context) (reportedAppreciations []ListReportedAppreciations, err error)
	GetReportedAppreciation(ctx context.Context, appreciationID int64) (reportedAppreciation ListReportedAppreciations, err error)
	DeleteAppreciation(ctx context.Context, moderationReq dto.ModerationReq) (err error)
	CheckResolution(ctx context.Context, id int64) (doesExist bool, appreciation_id int64, err error)
	ResolveAppreciation(ctx context.Context, moderationReq dto.ModerationReq) (err error)
	GetResolution(ctx context.Context, id int64) (reportedAppreciation ListReportedAppreciations, err error)
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
	Status           string `json:"status" db:"status"`
}

type ListReportedAppreciations struct {
	Id                int64          `db:"id"`
	Appreciation_id   int64          `db:"appreciation_id"`
	AppreciationDesc  string         `db:"appreciation_description"`
	TotalRewardPoints int64          `db:"total_reward_points"`
	Quarter           int64          `db:"quarter"`
	CoreValueName     string         `db:"core_value_name"`
	CoreValueDesc     string         `db:"core_value_description"`
	Sender            int64          `db:"sender"`
	Receiver          int64          `db:"receiver"`
	CreatedAt         int64          `db:"created_at"`
	IsValid           bool           `db:"is_valid"`
	ReportingComment  string         `db:"reporting_comment"`
	ReportedBy        int64          `db:"reported_by"`
	ReportedAt        int64          `db:"reported_at"`
	ModeratorComment  sql.NullString `db:"moderator_comment"`
	ModeratedBy       sql.NullInt64  `db:"moderated_by"`
	ModeratedAt       sql.NullInt64  `db:"moderated_at"`
	Status            string         `db:"status" json:"status"`
}
