package dto

import "github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"

type Appreciation struct {
	ID                int64  `json:"id"`
	CoreValueID       int    `json:"core_value_id" `
	Description       string `json:"description"`
	TotalRewardPoints int    `json:"total_reward_points,omitempty"`
	Quarter           int    `json:"quarter"`
	Sender            int64  `json:"sender"`
	Receiver          int64  `json:"receiver"`
	CreatedAt         int64  `json:"created_at"`
	UpdatedAt         int64  `json:"updated_at"`
}

type AppreciationFilter struct {
	Self      bool   `json:"Self"`
	Name      string `json:"sender_name"`
	SortOrder string `json:"sort_order"`
	Page      int16  `json:"page"`
	Limit     int16  `json:"page_size"`
}

type ResponseAppreciation struct {
	ID                  int    `json:"id"`
	CoreValueName       string `json:"core_value_name"`
	CoreValueDesc       string `json:"core_value_description"`
	Description         string `json:"description"`
	TotalRewardPoints   int    `json:"total_reward_points"`
	Quarter             string `json:"quarter"`
	SenderFirstName     string `json:"sender_first_name"`
	SenderLastName      string `json:"sender_last_name"`
	SenderImageURL      string `json:"sender_image_url"`
	SenderDesignation   string `json:"sender_designation"`
	ReceiverFirstName   string `json:"receiver_first_name"`
	ReceiverLastName    string `json:"receiver_last_name"`
	ReceiverImageURL    string `json:"receiver_image_url"`
	ReceiverDesignation string `json:"receiver_designation"`
	TotalRewards        int    `json:"total_rewards"`
	GivenRewardPoint    int    `json:"given_reward_point"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
}

// Pagination Object
type Pagination struct {
	CurrentPage int16 `json:"page"`
	TotalPage   int16 `json:"total_page"`
	PageSize    int16 `json:"page_size"`
}

type GetAppreciationResponse struct {
	Appreciations []ResponseAppreciation `json:"appreciations"`
	MetaData      Pagination             `json:"metadata"`
}

func (appr *Appreciation) ValidateCreateAppreciation() (err error) {

	if appr.CoreValueID <= 0 {
		return apperrors.InvalidCoreValueID
	}

	if appr.Description == "" {
		return apperrors.DescFieldBlank
	}

	if appr.Receiver <= 0 {
		return apperrors.InvalidReceiverID
	}

	return
}
