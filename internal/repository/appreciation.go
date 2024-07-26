package repository

import (
	"context"
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type AppreciationStorer interface {
	RepositoryTransaction

	CreateAppreciation(ctx context.Context, tx Transaction, appreciation dto.Appreciation) (Appreciation, error)
	GetAppreciationById(ctx context.Context, tx Transaction, appreciationId int32) (AppreciationInfo, error)
	GetAppreciations(ctx context.Context, tx Transaction, filter dto.AppreciationFilter) ([]AppreciationInfo, Pagination, error)
	DeleteAppreciation(ctx context.Context, tx Transaction,apprId int32) (bool, error)
	IsUserPresent(ctx context.Context, tx Transaction, userID int64) (bool, error)
	UpdateAppreciationTotalRewardsOfYesterday(ctx context.Context, tx Transaction) (bool, error)
	UpdateUserBadgesBasedOnTotalRewards(ctx context.Context, tx Transaction) (bool, error)
}

type Appreciation struct {
	ID                int64  `db:"id"`
	CoreValueID       int    `db:"core_value_id"`
	Description       string `db:"description"`
	IsValid           bool   `db:"is_valid"`
	TotalRewardPoints int    `db:"total_reward_points"`
	Quarter           int    `db:"quarter"`
	Sender            int64  `db:"sender"`
	Receiver          int64  `db:"receiver"`
	CreatedAt         int64  `db:"created_at"`
	UpdatedAt         int64  `db:"updated_at"`
}

type AppreciationInfo struct {
	ID                  int            `db:"id"`
	CoreValueName       string         `db:"core_value_name"`
	CoreValueDesc       string         `db:"core_value_description"`
	Description         string         `db:"description"`
	IsValid             bool           `db:"is_valid"`
	TotalRewardPoints   int            `db:"total_reward_points"`
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
	TotalRewards        int            `db:"total_rewards"`
	GivenRewardPoint    int            `db:"given_reward_point"`
	CreatedAt           int64          `db:"created_at"`
	UpdatedAt           int64          `db:"updated_at"`
}

// Pagination Object
type Pagination struct {
	RecordPerPage int16
	CurrentPage   int16
	TotalPage     int16
	TotalRecords  int32
}
