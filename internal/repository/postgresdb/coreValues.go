package repository

import (
	"context"

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
	//edit dele
	// deleteSubCoreValueQuery = `UPDATE core_values SET soft_delete = true, soft_delete_by = $1, updated_at = $2 WHERE org_id = $3 and parent_id = $4`
	// deleteCoreValueQuery    = `UPDATE core_values SET soft_delete = true, soft_delete_by = $1, updated_at = $2 WHERE org_id = $3 and id = $4`
	//edit dele
	updateCoreValueQuery = `UPDATE core_values SET (name, description) =
		($1, $2) where id = $3 RETURNING id, name, description, parent_core_value_id`

	// checkOrganisationQuery = `SELECT id from organizations WHERE id = $1`
	checkUniqueCoreVal = `SELECT id from core_values WHERE name = $1`
)

func (cs *coreValueStore) ListCoreValues(ctx context.Context) (coreValues []dto.ListCoreValuesResp, err error) {
	var DbResp []dto.ListCoreValuesRespDb
	err = cs.DB.SelectContext(
		ctx,
		&DbResp,
		listCoreValuesQuery,
	)

	for i := 0; i < len(DbResp); i++ {
		var coreValue dto.ListCoreValuesResp
		coreValue.ID = DbResp[i].ID
		coreValue.Name = DbResp[i].Name
		coreValue.Description = DbResp[i].Description
		coreValue.ParentCoreValueID = DbResp[i].ParentCoreValueID.Int64
		coreValues = append(coreValues, coreValue)
	}

	if err != nil {
		logger.WithFields(logger.Fields{
			"err": err.Error(),
		}).Error("Error while getting core values")
		return
	}

	return
}

func (cs *coreValueStore) GetCoreValue(ctx context.Context, coreValueID int64) (coreValue dto.GetCoreValueResp, err error) {
	err = cs.DB.GetContext(
		ctx,
		&coreValue,
		getCoreValueQuery,
		coreValueID,
	)
	if err != nil {
		logger.WithFields(logger.Fields{
			"err":         err.Error(),
			"coreValueId": coreValueID,
		}).Error("Error while getting core value")
		err = apperrors.InvalidCoreValueData
		return
	}

	return
}

func (cs *coreValueStore) CreateCoreValue(ctx context.Context, userId int64, coreValue dto.CreateCoreValueReq) (resp dto.CreateCoreValueResp, err error) {

	err = cs.DB.GetContext(
		ctx,
		&resp,
		createCoreValueQuery,
		coreValue.Name,
		coreValue.Description,
		coreValue.ParentCoreValueID,
	)
	if err != nil {
		logger.WithFields(logger.Fields{
			"err":               err.Error(),
			"core_value_params": coreValue,
		}).Error("Error while creating core value")
		return
	}

	return
}

func (cs *coreValueStore) UpdateCoreValue(ctx context.Context, coreValueID int64, updateReq dto.UpdateQueryRequest) (resp dto.UpdateCoreValuesResp, err error) {
	err = cs.DB.GetContext(
		ctx,
		&resp,
		updateCoreValueQuery,
		updateReq.Name,
		updateReq.Description,
		coreValueID,
	)
	if err != nil {
		logger.WithFields(logger.Fields{
			"err":           err.Error(),
			"core_value_id": coreValueID,
		}).Error("Error while updating core value")
		return
	}

	return
}

func (cs *coreValueStore) CheckUniqueCoreVal(ctx context.Context, text string) (isUnique bool, err error) {
	isUnique = false
	resp := []int64{}
	err = cs.DB.SelectContext(
		ctx,
		&resp,
		checkUniqueCoreVal,
		text,
	)

	if err != nil {
		logger.WithFields(logger.Fields{
			"err": err.Error(),
		}).Error("Error while checking unique core values")
		err = apperrors.InternalServerError
		return
	}

	if len(resp) <= 0 {
		isUnique = true
		return
	}

	return
}
