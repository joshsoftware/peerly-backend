package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	logger "github.com/sirupsen/logrus"
)

const (
	createOrganizationQuery = `INSERT INTO organization_config (
		id,
		reward_multiplier,
		reward_quota_renewal_frequency,
		timezone,
		created_by,updated_by)
		VALUES ($1, $2, $3, $4,$5,$6) RETURNING id,reward_multiplier,reward_quota_renewal_frequency,timezone,created_by,updated_by,created_at,updated_at`

	getOrganizationQuery = `SELECT id,
		reward_multiplier,
		reward_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at,updated_by FROM organization_config WHERE id=$1`
)

type OrganizationStore struct {
	BaseRepository
}

func NewOrganizationRepo(db *sqlx.DB) repository.OrganizationStorer {
	return &OrganizationStore{
		BaseRepository: BaseRepository{db}, // Use *sqlx.DB instead of *sql.DB
	}
}

func (s *OrganizationStore) CreateOrganizationConfig(ctx context.Context,tx repository.Transaction, org dto.OrganizationConfig) (createdOrganization repository.OrganizationConfig, err error) {

	insertQuery, args, err := sq.
	Insert("organization_config").Columns(constants.OrgConfigColumns...).
	Values(		1,
		org.RewardMultiplier,
		org.RewardQuotaRenewalFrequency,
		org.Timezone,
		org.CreatedBy,
		org.UpdatedBy,).
	PlaceholderFormat(sq.Dollar).
	Suffix("RETURNING \"id\",\"reward_multiplier\",\"reward_quota_renewal_frequency\",\"timezone\",\"created_by\",\"updated_by\",\"created_at\",\"updated_at\"").
	ToSql()
	if err != nil {
		logger.Error(err.Error())
		return repository.OrganizationConfig{}, apperrors.InternalServer
	}
	queryExecutor := s.InitiateQueryExecutor(tx)
	err = queryExecutor.QueryRowx(insertQuery, args...).Scan(&createdOrganization.ID,&createdOrganization.RewardMultiplier,&createdOrganization.RewardQuotaRenewalFrequency,&createdOrganization.Timezone,&createdOrganization.CreatedBy,&createdOrganization.UpdatedBy,&createdOrganization.CreatedAt,&createdOrganization.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error(apperrors.ErrRecordNotFound, "Just created an Organization, but can't find it!", err)
			return repository.OrganizationConfig{}, apperrors.InternalServer
		}
	}
	return
}

func (s *OrganizationStore) UpdateOrganizationConfig(ctx context.Context, tx repository.Transaction, reqOrganization dto.OrganizationConfig) (updatedOrganization repository.OrganizationConfig, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	updateBuilder := psql.Update("organization_config").
	Where(sq.Eq{"id": 1}).
	Suffix("RETURNING \"id\",\"reward_multiplier\",\"reward_quota_renewal_frequency\",\"timezone\",\"created_by\",\"updated_by\",\"created_at\",\"updated_at\"")

	if reqOrganization.RewardMultiplier != 0 {
		updateBuilder = updateBuilder.Set("reward_multiplier", reqOrganization.RewardMultiplier)
	}
	if reqOrganization.RewardQuotaRenewalFrequency != 0 {
		updateBuilder = updateBuilder.Set("reward_quota_renewal_frequency", reqOrganization.RewardQuotaRenewalFrequency)
	}
	if reqOrganization.Timezone != "" {
		updateBuilder = updateBuilder.Set("timezone", reqOrganization.Timezone)
	}

	updateBuilder = updateBuilder.Set("updated_at", time.Now().UnixMilli()).
		Set("updated_by", reqOrganization.UpdatedBy)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error building update query")
		return repository.OrganizationConfig{}, err
	}

	queryExecutor := s.InitiateQueryExecutor(tx)
	err = queryExecutor.QueryRowx(query, args...).Scan(&updatedOrganization.ID,&updatedOrganization.RewardMultiplier,&updatedOrganization.RewardQuotaRenewalFrequency,&updatedOrganization.Timezone,&updatedOrganization.CreatedBy,&updatedOrganization.UpdatedBy,&updatedOrganization.CreatedAt,&updatedOrganization.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error(apperrors.ErrRecordNotFound, "Just created an Organization, but can't find it!", err)
			return repository.OrganizationConfig{}, apperrors.InternalServer
		}
	}
	return
}

// GetOrganization - returns an organization from the database if it exists based on its ID primary key
func (s *OrganizationStore) GetOrganizationConfig(ctx context.Context, tx repository.Transaction) (organization repository.OrganizationConfig, err error) {
	queryExecutor := s.InitiateQueryExecutor(tx)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	queryBuilder := psql.
	Select("*").
	From("organization_config").
	Where(sq.Eq{"id": 1})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error building select query")
		return repository.OrganizationConfig{}, err
	}

	err = queryExecutor.QueryRowx( query, args...).Scan(&organization.ID,&organization.RewardMultiplier,&organization.RewardQuotaRenewalFrequency,&organization.Timezone,&organization.CreatedAt,&organization.CreatedBy,&organization.UpdatedAt,&organization.UpdatedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("organizationID", 1).Warn("Organization not found")
			return repository.OrganizationConfig{}, apperrors.OrganizationConfigNotFound
		}
		logger.WithField("err", err.Error()).Error("Error fetching organization")
		return repository.OrganizationConfig{}, err
	}

	return organization, nil
}
///helper functions Organization

func OrganizationConfigToDB(org dto.OrganizationConfig) repository.OrganizationConfig {
	return repository.OrganizationConfig{
		RewardMultiplier:            org.RewardMultiplier,
		ID:                          org.ID,
		RewardQuotaRenewalFrequency: org.RewardQuotaRenewalFrequency,
		Timezone:                    org.Timezone,
		CreatedAt:                   org.CreatedAt,
		CreatedBy:                   org.CreatedBy,
		UpdatedAt:                   org.UpdatedAt,
		UpdatedBy:                   org.UpdatedBy,
	}
}
