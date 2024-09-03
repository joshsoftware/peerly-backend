package repository

import (
	"context"
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type GradesStorer interface {
	ListGrades(ctx context.Context) (gradesList []Grade, err error)
	EditGrade(ctx context.Context, reqData dto.UpdateGradeReq) (err error)
}

type Grade struct {
	Id        int64         `db:"id"`
	Name      string        `db:"name"`
	Points    int64         `db:"points"`
	UpdatedBy sql.NullInt64 `db:"updated_by"`
}
