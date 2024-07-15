package repository

import (
	"context"
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type AppreciationStorer interface {
	RepositoryTransaction

	CreateAppreciation(ctx context.Context, tx Transaction, appreciation dto.Appreciation) (Appreciation, error)
	GetAppreciationById(ctx context.Context, tx Transaction, appreciationId int) (AppreciationInfo, error)
	GetAppreciation(ctx context.Context, tx Transaction, filter dto.AppreciationFilter) ([]AppreciationInfo, Pagination, error)
	ValidateAppreciation(ctx context.Context, tx Transaction, isValid bool, apprId int) (bool, error)
	IsUserPresent(ctx context.Context, tx Transaction, userID int64) (bool, error)
	UpdateAppreciationTotalRewardsOfYesterday(ctx context.Context, tx Transaction) (bool, error)
	UpdateUserBadgesBasedOnTotalRewards(ctx context.Context, tx Transaction) (bool, error)
}

type Appreciation struct {
	ID           int64  `db:"id"`
	CoreValueID  int    `db:"core_value_id"`
	Description  string `db:"description"`
	IsValid      bool   `db:"is_valid"`
	TotalRewards int    `db:"total_rewards"`
	Quarter      int    `db:"quarter"`
	Sender       int64  `db:"sender"`
	Receiver     int64  `db:"receiver"`
	CreatedAt    int64  `db:"created_at"`
	UpdatedAt    int64  `db:"updated_at"`
}

type AppreciationInfo struct {
	ID                  int            `db:"id"`
	CoreValueName       string         `db:"core_value_name"`
	Description         string         `db:"description"`
	IsValid             bool           `db:"is_valid"`
	TotalRewards        int            `db:"total_reward_points"`
	Quarter             string         `db:"quarter"`
	SenderId            int64          `db:"sender_id"`
	SenderFirstName     string         `db:"sender_first_name"`
	SenderLastName      string         `db:"sender_last_name"`
	SenderImageURL      sql.NullString `db:"sender_image_url"`
	SenderDesignation   string         `db:"sender_designation"`
	ReceiverId          int64          `db:"receiver_id"`
	ReceiverFirstName   string         `db:"receiver_first_name"`
	ReceiverLastName    string         `db:"receiver_last_name"`
	ReceiverImageURL    sql.NullString `db:"receiver_image_url"`
	ReceiverDesignation string         `db:"receiver_designation"`
	CreatedAt           int64          `db:"created_at"`
	UpdatedAt           int64          `db:"updated_at"`
}

// Pagination Object
type Pagination struct {
	Next          *int64
	Previous      *int64
	RecordPerPage int64
	CurrentPage   int64
	TotalPage     int64
	TotalRecords  int64
}
