package dto

type ChangeLogLevelRequest struct {
	LogLevel     string `json:"loglevel"`
	DeveloperKey string `json:"developer_key"`
}
