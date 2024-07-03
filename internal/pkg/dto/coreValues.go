package dto

type UpdateQueryRequest struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type UpdateCoreValuesResp struct {
	ID                int64  `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Description       string `json:"description" db:"description"`
	ParentCoreValueID *int64 `json:"parent_core_value_id" db:"parent_core_value_id"`
}

type ListCoreValuesResp struct {
	ID                int64  `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Description       string `json:"description" db:"description"`
	ParentCoreValueID *int64 `json:"parent_core_value_id" db:"parent_core_value_id"`
}

type GetCoreValueResp struct {
	ID                int64  `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Description       string `json:"description" db:"description"`
	ParentCoreValueID *int64 `json:"parent_core_value_id" db:"parent_core_value_id"`
}

type CreateCoreValueReq struct {
	Name              string `json:"name" db:"name"`
	Description       string `json:"description" db:"description"`
	ParentCoreValueID *int64 `json:"parent_core_value_id" db:"parent_core_value_id"`
}

type CreateCoreValueResp struct {
	ID                int64  `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Description       string `json:"description" db:"description"`
	ParentCoreValueID *int64 `json:"parent_core_value_id" db:"parent_core_value_id"`
}
