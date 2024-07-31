package dto

type ReportAppreciationReq struct {
	AppreciationId   int64  `json:"appreciation_id" db:"appreciation_id"`
	ReportingComment string `json:"reporting_comment" db:"reporting_comment"`
	ReportedBy       int64  `json:"reported_by" db:"reported_by"`
}

type ReportAppricaitionResp struct {
	Id               int64  `json:"id" db:"id"`
	AppreciationId   int64  `json:"appreciation_id" db:"appreciation_id"`
	ReportingComment string `json:"reporting_comment" db:"reporting_comment"`
	ReportedBy       int64  `json:"reported_by" db:"reported_by"`
	ReportedAt       int64  `json:"reported_at" db:"reported_at"`
}

type GetSenderAndReceiverResp struct {
	Sender   int64 `json:"sender" db:"sender"`
	Receiver int64 `json:"receiver" db:"receiver"`
}

type ReportedAppreciation struct {
	Id                  int64  `json:"id"`
	Appreciation_id     int64  `json:"appreciation_id"`
	AppreciationDesc    string `json:"appreciation_description"`
	TotalRewardPoints   int64  `json:"total_reward_points"`
	Quarter             int64  `json:"quarter"`
	CoreValueName       string `json:"core_value_name"`
	CoreValueDesc       string `json:"core_value_description"`
	SenderFirstName     string `json:"sender_first_name"`
	SenderLastName      string `json:"sender_last_name"`
	SenderImgUrl        string `json:"sender_image_url"`
	SenderDesignation   string `json:"sender_designation"`
	ReceiverFirstName   string `json:"receiver_first_name"`
	ReceiverLastName    string `json:"receiver_last_name"`
	ReceiverImgUrl      string `json:"receiver_image_url"`
	ReceiverDesignation string `json:"receiver_designation"`
	ReportingComment    string `json:"reporting_comment"`
	ReportedBy          int64  `json:"reported_by"`
	ReportedAt          int64  `json:"reported_at"`
}

type ModerationReq struct {
	ResolutionId     int64
	AppreciationId   int64
	ModeratorComment string `json:"moderator_comment"`
	ModeratedBy      int64
}
