package repository

import (
	"context"
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type CoreValueStorer interface {
	ListCoreValues(ctx context.Context) (coreValues []CoreValue, err error)
	GetCoreValue(ctx context.Context, coreValueID int64) (coreValue CoreValue, err error)
	CreateCoreValue(ctx context.Context, coreValue dto.CreateCoreValueReq) (resp CoreValue, err error)
	UpdateCoreValue(ctx context.Context, coreValue dto.UpdateQueryRequest) (resp CoreValue, err error)
	// CheckOrganisation(ctx context.Context, organisationId int64) (err error)
	CheckUniqueCoreVal(ctx context.Context, text string) (res bool, err error)
}

// CoreValue - struct representing a core value object
type CoreValue struct {
	ID                int64         `db:"id"`
	Name              string        `db:"name"`
	Description       string        `db:"description"`
	ParentCoreValueID sql.NullInt64 `db:"parent_core_value_id"`
}
