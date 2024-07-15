package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

var (
	sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

var CoreValueColumns = []string{"id", "name", "description", "parent_core_value_id"}

type coreValueStore struct {
	DB        *sqlx.DB
	TableName string
}

func NewCoreValueRepo(db *sqlx.DB) repository.CoreValueStorer {
	return &coreValueStore{
		DB:        db,
		TableName: "core_values",
	}
}

func (cs *coreValueStore) ListCoreValues(ctx context.Context) (coreValues []repository.CoreValue, err error) {
	queryBuilder := sq.Select(CoreValueColumns...).From(cs.TableName)
	listCoreValuesQuery, _, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
		return
	}
	err = cs.DB.SelectContext(
		ctx,
		&coreValues,
		listCoreValuesQuery,
	)

	if err != nil {
		err = fmt.Errorf("error while getting core values, err: %w", err)
		return
	}

	return
}

func (cs *coreValueStore) GetCoreValue(ctx context.Context, coreValueID int64) (coreValue repository.CoreValue, err error) {
	queryBuilder := sq.
		Select(CoreValueColumns...).
		From(cs.TableName).
		Where(squirrel.Eq{"id": coreValueID})

	getCoreValueQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Errorf("error in generating squirrel query, err: %w", err)
		err = apperrors.InternalServerError
		return
	}

	err = cs.DB.GetContext(
		ctx,
		&coreValue,
		getCoreValueQuery,
		args...,
	)
	if err != nil {
		logger.Errorf("error while getting core value, corevalue_id: %d, err: %w", coreValueID, err)
		err = apperrors.InvalidCoreValueData
		return
	}

	return
}

func (cs *coreValueStore) CreateCoreValue(ctx context.Context, coreValue dto.CreateCoreValueReq) (resp repository.CoreValue, err error) {

	queryBuilder := sq.Insert(cs.TableName).Columns(CoreValueColumns[1:]...).Values(coreValue.Name, coreValue.Description, coreValue.ParentCoreValueID).Suffix("RETURNING id, name, description, parent_core_value_id")

	createCoreValueQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
		return
	}

	err = cs.DB.GetContext(
		ctx,
		&resp,
		createCoreValueQuery,
		args...,
	)
	if err != nil {
		err = fmt.Errorf("error while creating core value, err: %w", err)
		return
	}

	return
}

func (cs *coreValueStore) UpdateCoreValue(ctx context.Context, updateReq dto.UpdateQueryRequest) (resp repository.CoreValue, err error) {
	queryBuilder := sq.Update(cs.TableName).
		Set("name", updateReq.Name).
		Set("description", updateReq.Description).
		Where(squirrel.Eq{"id": updateReq.Id}).
		Suffix("RETURNING id, name, description, parent_core_value_id")

	updateCoreValueQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
		return
	}
	err = cs.DB.GetContext(
		ctx,
		&resp,
		updateCoreValueQuery,
		args...,
	)
	if err != nil {
		err = fmt.Errorf("error while updating core value, corevalue_id: %d, err: %w", updateReq.Id, err)
		return
	}

	return
}

func (cs *coreValueStore) CheckUniqueCoreVal(ctx context.Context, name string) (isUnique bool, err error) {

	isUnique = false
	resp := []int64{}
	queryBuilder := sq.Select("id").
		From(cs.TableName).
		Where(squirrel.Like{"name": name})

	checkUniqueCoreVal, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
		return
	}

	err = cs.DB.SelectContext(
		ctx,
		&resp,
		checkUniqueCoreVal,
		args...,
	)

	if err != nil {
		err = fmt.Errorf("error while checking unique core vlaues, err: %w", err)
		return
	}

	if len(resp) <= 0 {
		isUnique = true
		return
	}

	return
}
