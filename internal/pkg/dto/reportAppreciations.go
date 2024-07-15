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
