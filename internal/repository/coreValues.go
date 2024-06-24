package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

type coreValueStore struct {
	DB *sqlx.DB
}

type CoreValueStorer interface {
	ListCoreValues(ctx context.Context, organisationID int64) (coreValues []dto.ListCoreValuesResp, err error)
	GetCoreValue(ctx context.Context, organisationID, coreValueID int64) (coreValue dto.GetCoreValueResp, err error)
	CreateCoreValue(ctx context.Context, organisationID int64, userId int64, coreValue dto.CreateCoreValueReq) (resp dto.CreateCoreValueResp, err error)
	DeleteCoreValue(ctx context.Context, organisationID, coreValueID int64, userId int64) (err error)
	UpdateCoreValue(ctx context.Context, organisationID, coreValueID int64, coreValue dto.UpdateQueryRequest) (resp dto.UpdateCoreValuesResp, err error)
	CheckOrganisation(ctx context.Context, organisationId int64) (err error)
	CheckUniqueCoreVal(ctx context.Context, organisationId int64, text string) (res bool, err error)
}

func NewCoreValueRepo(db *sqlx.DB) CoreValueStorer {
	return &coreValueStore{
		DB: db,
	}
}

const (
	listCoreValuesQuery  = `SELECT id, org_id, text, description, parent_id, created_at, updated_at  FROM core_values WHERE org_id = $1 and soft_delete = false`
	getCoreValueQuery    = `SELECT id, org_id, text, description, parent_id, soft_delete FROM core_values WHERE org_id = $1 and id = $2`
	createCoreValueQuery = `INSERT INTO core_values (org_id, text,
		description, parent_id, thumbnail_url,created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, org_id, text, description, parent_id, thumbnail_url, created_by, created_at`
	// deleteSubCoreValueQuery = `DELETE FROM core_values WHERE orgId = $1 and parentId = $2`
	deleteSubCoreValueQuery = `UPDATE core_values SET soft_delete = true, soft_delete_by = $1, updated_at = $2 WHERE org_id = $3 and parent_id = $4`
	// deleteCoreValueQuery    = `DELETE FROM core_values WHERE orgId = $1 and id = $2`
	deleteCoreValueQuery = `UPDATE core_values SET soft_delete = true, soft_delete_by = $1, updated_at = $2 WHERE org_id = $3 and id = $4`
	updateCoreValueQuery = `UPDATE core_values SET (text, description, thumbnail_url, updated_at) =
		($1, $2, $3, $4) where id = $5 and org_id = $6 RETURNING id, org_id, text, description, parent_id, thumbnail_url, updated_at`

	checkOrganisationQuery = `SELECT id from organizations WHERE id = $1`
	checkUniqueCoreVal     = `SELECT id from core_values WHERE org_id = $1 and text = $2`
)

// CoreValue - struct representing a core value object
type CoreValue struct {
	ID           int64     `db:"id" json:"id"`
	OrgID        int64     `db:"org_id" json:"org_id"`
	Text         string    `db:"text" json:"text"`
	Description  string    `db:"description" json:"description"`
	ParentID     *int64    `db:"parent_id" json:"parent_id"`
	ThumbnailURL *string   `db:"thumbnail_url" json:"thumbnail_url"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	SoftDelete   bool      `db:"soft_delete"`
	SoftDeleteBy int64     `db:"soft_delete_by"`
	CreatedBy    int64     `db:"created_by"`
}

func (cs *coreValueStore) ListCoreValues(ctx context.Context, organisationID int64) (coreValues []dto.ListCoreValuesResp, err error) {
	err = cs.DB.SelectContext(
		ctx,
		&coreValues,
		listCoreValuesQuery,
		organisationID,
	)

	if err != nil {
		logger.WithFields(logger.Fields{
			"err":   err.Error(),
			"orgId": organisationID,
		}).Error("Error while getting core values")
		return
	}

	return
}

func (cs *coreValueStore) GetCoreValue(ctx context.Context, organisationID, coreValueID int64) (coreValue dto.GetCoreValueResp, err error) {
	err = cs.DB.GetContext(
		ctx,
		&coreValue,
		getCoreValueQuery,
		organisationID,
		coreValueID,
	)
	if err != nil {
		logger.WithFields(logger.Fields{
			"err":         err.Error(),
			"orgId":       organisationID,
			"coreValueId": coreValueID,
		}).Error("Error while getting core value")
		err = apperrors.InvalidCoreValueData
		return
	}

	return
}

func (cs *coreValueStore) CreateCoreValue(ctx context.Context, organisationID int64, userId int64, coreValue dto.CreateCoreValueReq) (resp dto.CreateCoreValueResp, err error) {
	now := time.Now()
	err = cs.DB.GetContext(
		ctx,
		&resp,
		createCoreValueQuery,
		organisationID,
		coreValue.Text,
		coreValue.Description,
		coreValue.ParentID,
		coreValue.ThumbnailURL,
		userId,
		now,
		now,
	)
	if err != nil {
		logger.WithFields(logger.Fields{
			"err":               err.Error(),
			"org_id":            organisationID,
			"core_value_params": coreValue,
		}).Error("Error while creating core value")
		return
	}

	return
}

func (cs *coreValueStore) DeleteCoreValue(ctx context.Context, organisationID, coreValueID int64, userId int64) (err error) {
	now := time.Now()
	_, err = cs.DB.ExecContext(
		ctx,
		deleteSubCoreValueQuery,
		userId,
		now,
		organisationID,
		coreValueID,
	)
	if err != nil {
		logger.WithFields(logger.Fields{
			"err":         err.Error(),
			"orgId":       organisationID,
			"coreValueId": coreValueID,
		}).Error("Error while deleting sub core value")
		return
	}

	_, err = cs.DB.ExecContext(
		ctx,
		deleteCoreValueQuery,
		userId,
		now,
		organisationID,
		coreValueID,
	)
	if err != nil {
		logger.WithFields(logger.Fields{
			"err":           err.Error(),
			"org_id":        organisationID,
			"core_value_id": coreValueID,
		}).Error("Error while deleting core value")
		return
	}

	return
}

func (cs *coreValueStore) UpdateCoreValue(ctx context.Context, organisationID int64, coreValueID int64, updateReq dto.UpdateQueryRequest) (resp dto.UpdateCoreValuesResp, err error) {
	now := time.Now()
	err = cs.DB.GetContext(
		ctx,
		&resp,
		updateCoreValueQuery,
		updateReq.Text,
		updateReq.Description,
		updateReq.ThumbnailUrl,
		now,
		coreValueID,
		organisationID,
	)
	if err != nil {
		logger.WithFields(logger.Fields{
			"err":           err.Error(),
			"org_id":        organisationID,
			"core_value_id": coreValueID,
		}).Error("Error while updating core value")
		return
	}

	return
}

func (cs *coreValueStore) CheckOrganisation(ctx context.Context, organisationId int64) (err error) {
	resp := []int64{}
	err = cs.DB.SelectContext(
		ctx,
		&resp,
		checkOrganisationQuery,
		organisationId,
	)

	if len(resp) <= 0 {
		err = apperrors.InvalidOrgId
	}

	if err != nil {
		logger.WithFields(logger.Fields{
			"err":    err.Error(),
			"org_id": organisationId,
		}).Error("Error while checking organisation")
		err = apperrors.InvalidOrgId
	}

	return
}

func (cs *coreValueStore) CheckUniqueCoreVal(ctx context.Context, organisationId int64, text string) (isUnique bool, err error) {
	isUnique = false
	resp := []int64{}
	err = cs.DB.SelectContext(
		ctx,
		&resp,
		checkUniqueCoreVal,
		organisationId,
		text,
	)

	fmt.Println("resp: ", resp)
	fmt.Println("err: ", err)

	if err != nil {
		logger.WithFields(logger.Fields{
			"err": err.Error(),
		}).Error("Error while checking organisation")
		err = apperrors.InternalServerError
		return
	}

	if len(resp) <= 0 {
		isUnique = true
		return
	}

	return
}
