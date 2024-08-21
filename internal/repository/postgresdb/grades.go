package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

var (
	gradeColumns = []string{"id", "name", "points", "updated_by"}
)

type gradeStore struct {
	DB          *sqlx.DB
	GradesTable string
}

func NewGradesRepo(db *sqlx.DB) repository.GradesStorer {
	return &gradeStore{
		DB:          db,
		GradesTable: constants.GradesTable,
	}
}

func (gs *gradeStore) ListGrades(ctx context.Context) (gradesList []repository.Grade, err error) {

	queryBuilder := repository.Sq.Select(gradeColumns...).From(gs.GradesTable).OrderBy("id")
	getGradesQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
		return
	}

	err = gs.DB.SelectContext(ctx, &gradesList, getGradesQuery, args...)
	if err != nil {
		err = fmt.Errorf("error fetching grades from database, err: %w", err)
		return
	}

	return
}

func (gs *gradeStore) EditGrade(ctx context.Context, reqData dto.UpdateGradeReq) (err error) {
	queryBuilder := repository.Sq.Update(gs.GradesTable).Set("points", reqData.Points).Set("updated_by", reqData.UpdatedBy).Where(squirrel.Eq{"id": reqData.Id})
	updateGradeQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
		return
	}
	_, err = gs.DB.ExecContext(ctx, updateGradeQuery, args...)
	if err != nil {
		err = fmt.Errorf("error in updating grade points, err: %w", err)
		return
	}

	return
}
