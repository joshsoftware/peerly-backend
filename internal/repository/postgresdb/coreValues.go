package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type coreValueStore struct {
	DB *sqlx.DB
}

func NewCoreValueRepo(db *sqlx.DB) repository.CoreValueStorer {
	return &coreValueStore{
		DB: db,
	}
}

const (
	listCoreValuesQuery  = `SELECT id, name, description, parent_core_value_id FROM core_values`
	getCoreValueQuery    = `SELECT id, name, description, parent_core_value_id FROM core_values WHERE id = $1`
	createCoreValueQuery = `INSERT INTO core_values (name,
		description, parent_core_value_id) VALUES ($1, $2, $3) RETURNING id, name, description, parent_core_value_id`
	updateCoreValueQuery = `UPDATE core_values SET (name, description) =
		($1, $2) where id = $3 RETURNING id, name, description, parent_core_value_id`

	checkUniqueCoreVal = `SELECT id from core_values WHERE name = $1`
)

func (cs *coreValueStore) ListCoreValues(ctx context.Context) (coreValues []repository.CoreValue, err error) {
	err = cs.DB.SelectContext(
		ctx,
		&coreValues,
		listCoreValuesQuery,
	)

	if err != nil {
		logger.Error(fmt.Sprintf("error while getting core values, err: %s", err.Error()))
		return
	}

	return
}

func (cs *coreValueStore) GetCoreValue(ctx context.Context, coreValueID int64) (coreValue repository.CoreValue, err error) {
	err = cs.DB.GetContext(
		ctx,
		&coreValue,
		getCoreValueQuery,
		coreValueID,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("error while getting core value, corevalue_id: %d, err: %s", coreValueID, err.Error()))
		err = apperrors.InvalidCoreValueData
		return
	}

	return
}

func (cs *coreValueStore) CreateCoreValue(ctx context.Context, coreValue dto.CreateCoreValueReq) (resp repository.CoreValue, err error) {

	err = cs.DB.GetContext(
		ctx,
		&resp,
		createCoreValueQuery,
		coreValue.Name,
		coreValue.Description,
		coreValue.ParentCoreValueID,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("error while creating core value, err: %s", err.Error()))
		return
	}

	return
}

func (cs *coreValueStore) UpdateCoreValue(ctx context.Context, updateReq dto.UpdateQueryRequest) (resp repository.CoreValue, err error) {
	err = cs.DB.GetContext(
		ctx,
		&resp,
		updateCoreValueQuery,
		updateReq.Name,
		updateReq.Description,
		updateReq.Id,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("error while updating core value, corevalue_id: %d, err: %s", updateReq.Id, err.Error()))
		return
	}

	return
}

func (cs *coreValueStore) CheckUniqueCoreVal(ctx context.Context, name string) (isUnique bool, err error) {
	isUnique = false
	resp := []int64{}
	err = cs.DB.SelectContext(
		ctx,
		&resp,
		checkUniqueCoreVal,
		name,
	)

	if err != nil {
		logger.Error(fmt.Sprintf("error while checking unique core vlaues, err: %s", err.Error()))
		err = apperrors.InternalServerError
		return
	}

	if len(resp) <= 0 {
		isUnique = true
		return
	}

	return
}
