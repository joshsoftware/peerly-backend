package dto

type CoreValue struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ParentCoreValueID int64  `json:"parent_id"`
}

type UpdateQueryRequest struct {
	Id          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type CreateCoreValueReq struct {
	Name              string `json:"name" db:"name"`
	Description       string `json:"description" db:"description"`
	ParentCoreValueID *int64 `json:"parent_core_value_id" db:"parent_core_value_id"`
}
