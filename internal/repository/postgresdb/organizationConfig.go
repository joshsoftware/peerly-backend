package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type OrganizationConfigStore struct {
	BaseRepository
	OrganizationConfigTable string
}

func NewOrganizationConfigRepo(db *sqlx.DB) repository.OrganizationConfigStorer {
	return &OrganizationConfigStore{
		BaseRepository:          BaseRepository{db},
		OrganizationConfigTable: constants.OrganizationConfigTable,
	}
}

func (org *OrganizationConfigStore) CreateOrganizationConfig(ctx context.Context, tx repository.Transaction, orgConfigInfo dto.OrganizationConfig) (createdOrganization repository.OrganizationConfig, err error) {

	logger.Infof(ctx, "orgRepo: CreateOrganizationConfig: %v", orgConfigInfo)
	queryExecutor := org.InitiateQueryExecutor(tx)

	insertQuery, args, err := repository.Sq.
		Insert(org.OrganizationConfigTable).Columns(constants.OrgConfigColumns...).
		Values(1,
			orgConfigInfo.RewardMultiplier,
			orgConfigInfo.RewardQuotaRenewalFrequency,
			orgConfigInfo.Timezone,
			orgConfigInfo.CreatedBy,
			orgConfigInfo.UpdatedBy).
		Suffix("RETURNING id, reward_multiplier ,reward_quota_renewal_frequency, timezone, created_by, updated_by, created_at, updated_at").
		ToSql()
	if err != nil {
		logger.Errorf(ctx, "err in creating query: %v", err)
		return repository.OrganizationConfig{}, apperrors.InternalServer
	}

	logger.Debug(ctx, fmt.Sprintf("orgRepo: query: %s,args: %v", insertQuery, args))

	err = queryExecutor.QueryRowx(insertQuery, args...).StructScan(&createdOrganization)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf(ctx, "err in creating orgconfig : %v ", err)
			return repository.OrganizationConfig{}, apperrors.InternalServer
		}
	}

	logger.Debug(ctx, fmt.Sprintf("orgRepo: query: %s,args: %v", insertQuery, args))
	return
}

func (org *OrganizationConfigStore) UpdateOrganizationConfig(ctx context.Context, tx repository.Transaction, reqOrganization dto.OrganizationConfig) (updatedOrganization repository.OrganizationConfig, err error) {
	logger.Infof(ctx, "orgRepo: UpdateOrganizationConfig: %v", reqOrganization)
	queryExecutor := org.InitiateQueryExecutor(tx)

	updateBuilder := repository.Sq.Update(org.OrganizationConfigTable).
		Where(sq.Eq{"id": constants.DefaultOrgID}).
		Suffix("RETURNING id, reward_multiplier, reward_quota_renewal_frequency, timezone, created_by, updated_by, created_at, updated_at")

	if reqOrganization.RewardMultiplier != 0 {
		updateBuilder = updateBuilder.Set("reward_multiplier", reqOrganization.RewardMultiplier)
	}
	if reqOrganization.RewardQuotaRenewalFrequency != 0 {
		updateBuilder = updateBuilder.Set("reward_quota_renewal_frequency", reqOrganization.RewardQuotaRenewalFrequency)
	}
	if reqOrganization.Timezone != "" {
		updateBuilder = updateBuilder.Set("timezone", reqOrganization.Timezone)
	}

	updateBuilder = updateBuilder.
		Set("updated_at", time.Now().UnixMilli()).
		Set("updated_by", reqOrganization.UpdatedBy)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		logger.Errorf(ctx, "Error building update query: %v", err)
		return repository.OrganizationConfig{}, err
	}

	logger.Debug(ctx, fmt.Sprintf("orgRepo: query: %s,args: %v", query, args))
	err = queryExecutor.QueryRowx(query, args...).StructScan(&updatedOrganization)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf(ctx, "orgRepo: err in updating orgconfig : %v ", err)
			return repository.OrganizationConfig{}, apperrors.InternalServer
		}
	}
	return
}

// GetOrganization - returns an organization from the database if it exists based on its ID primary key
func (org *OrganizationConfigStore) GetOrganizationConfig(ctx context.Context, tx repository.Transaction) (updatedOrgConfig repository.OrganizationConfig, err error) {
	logger.Debug(ctx, "orgRepo: GetOrganizationConfig")
	queryExecutor := org.InitiateQueryExecutor(tx)

	queryBuilder := repository.Sq.
		Select(constants.OrgConfigColumns...).
		From(org.OrganizationConfigTable).
		Where(sq.Eq{"id": constants.DefaultOrgID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Errorf(ctx, "Error building select query: %v", err)
		return repository.OrganizationConfig{}, err
	}

	logger.Debug(ctx, fmt.Sprintf("orgRepo: query: %s,args: %v", query, args))

	err = queryExecutor.QueryRowx(query, args...).StructScan(&updatedOrgConfig)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf(ctx, "Organization not found: %v", err)
			return repository.OrganizationConfig{}, apperrors.OrganizationConfigNotFound
		}
		logger.Errorf(ctx, "Error fetching organization: %v", err)
		return repository.OrganizationConfig{}, err
	}

	logger.Debug(ctx, " updateOrgConfig: ", updatedOrgConfig)
	return updatedOrgConfig, nil
}
