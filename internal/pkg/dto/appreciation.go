package dto

type Appreciation struct {
	ID           int64  `json:"id"`
	CoreValueID  int    `json:"core_value_id" `
	Description  string `json:"description"`
	TotalRewards int    `json:"total_rewards,omitempty"`
	Quarter      int    `json:"quarter"`
	Sender       int64  `json:"sender"`
	Receiver     int64  `json:"receiver"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type AppreciationFilter struct {
	Name      string `json:"sender_name"`
	SortOrder string `json:"sort_order"`
}

type ResponseAppreciation struct {
	ID                  int    `json:"id"`
	CoreValueName       string `json:"core_value_name"`
	Description         string `json:"description"`
	IsValid             bool   `json:"is_valid"`
	TotalRewards        int    `json:"total_rewards"`
	Quarter             string `json:"quarter"`
	SenderFirstName     string `json:"sender_first_name"`
	SenderLastName      string `json:"sender_last_name"`
	SenderImageURL      string `json:"sender_image_url"`
	SenderDesignation   string `json:"sender_designation"`
	ReceiverFirstName   string `json:"receiver_first_name"`
	ReceiverLastName    string `json:"receiver_last_name"`
	ReceiverImageURL    string `json:"receiver_image_url"`
	ReceiverDesignation string `json:"receiver_designation"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
}

func (appr *Appreciation)CreateAppreciation() (errorResponse ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)

	if appr.CoreValueID <= 0 {
		fieldErrors["core_value_id"] = "enter valid core value id"
	}

	if appr.Description == "" {
		fieldErrors["description"] = "enter description"
	}

	if appr.Receiver <= 0 {
		fieldErrors["receiver"] = "enter valid receiver id"
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	errorResponse = ErrorResponse{
		Error: ErrorObject{
			Code:          "invalid_data",
			MessageObject: MessageObject{Message: "Please provide valid appreciation data"},
			Fields:        fieldErrors,
		},
	}

	return
}
