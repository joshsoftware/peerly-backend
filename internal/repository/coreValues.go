package repository

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type CoreValueStorer interface {
	ListCoreValues(ctx context.Context) (coreValues []dto.ListCoreValuesResp, err error)
	GetCoreValue(ctx context.Context, coreValueID int64) (coreValue dto.GetCoreValueResp, err error)
	CreateCoreValue(ctx context.Context, userId int64, coreValue dto.CreateCoreValueReq) (resp dto.CreateCoreValueResp, err error)
	UpdateCoreValue(ctx context.Context, coreValueID int64, coreValue dto.UpdateQueryRequest) (resp dto.UpdateCoreValuesResp, err error)
	// CheckOrganisation(ctx context.Context, organisationId int64) (err error)
	CheckUniqueCoreVal(ctx context.Context, text string) (res bool, err error)
}

// CoreValue - struct representing a core value object
type CoreValue struct {
	ID                int64  `db:"id" json:"id"`
	Name              string `db:"name" json:"name"`
	Description       string `db:"description" json:"description"`
	ParentCoreValueID *int64 `db:"parent_core_value_id" json:"parent_core_value_id"`
}
