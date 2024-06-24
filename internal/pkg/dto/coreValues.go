package dto

import "time"

type UpdateQueryRequest struct {
	Text         string `db:"text" json:"text"`
	Description  string `db:"description" json:"description"`
	ThumbnailUrl string `db:"thumbnailurl" json:"thumbnail_url"`
}

type UpdateCoreValuesResp struct {
	ID           int64     `json:"id" db:"id"`
	OrgID        int64     `json:"org_id" db:"org_id"`
	Text         string    `json:"text" db:"text"`
	Description  string    `json:"description" db:"description"`
	ParentID     *int64    `json:"parent_id" db:"parent_id"`
	ThumbnailURL *string   `json:"thumbnail_url" db:"thumbnail_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    int64     `json:"created_by" db:"created_by"`
}

type ListCoreValuesResp struct {
	ID           int64     `json:"id" db:"id"`
	OrgID        int64     `json:"org_id" db:"org_id"`
	Text         string    `json:"text" db:"text"`
	Description  string    `json:"description" db:"description"`
	ParentID     *int64    `json:"parent_id" db:"parent_id"`
	ThumbnailURL *string   `json:"thumbnail_url" db:"thumbnail_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    int64     `json:"created_by" db:"created_by"`
}

type GetCoreValueResp struct {
	ID           int64     `json:"id" db:"id"`
	OrgID        int64     `json:"org_id" db:"org_id"`
	Text         string    `json:"text" db:"text"`
	Description  string    `json:"description" db:"description"`
	ParentID     *int64    `json:"parent_id" db:"parent_id"`
	ThumbnailURL *string   `json:"thumbnail_url" db:"thumbnail_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    int64     `json:"created_by" db:"created_by"`
	SoftDelete   bool      `json:"soft_delete" db:"soft_delete"`
	SoftDeleteBy int64     `json:"soft_delete_by" db:"soft_delete_by"`
}

type CreateCoreValueReq struct {
	Text         string `json:"text" db:"text"`
	Description  string `json:"description" db:"description"`
	ParentID     *int64 `json:"parent_id" db:"parent_id"`
	ThumbnailURL string `json:"thumbnail_url" db:"thumbnail_url"`
}

type CreateCoreValueResp struct {
	ID           int64     `json:"id" db:"id"`
	OrgID        int64     `json:"org_id" db:"org_id"`
	Text         string    `json:"text" db:"text"`
	Description  string    `json:"description" db:"description"`
	ParentID     *int64    `json:"parent_id" db:"parent_id"`
	ThumbnailURL *string   `json:"thumbnail_url" db:"thumbnail_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    int64     `json:"created_by" db:"created_by"`
	SoftDelete   bool      `json:"soft_delete" db:"soft_delete"`
	SoftDeleteBy int64     `json:"soft_delete_by" db:"soft_delete_by"`
}
