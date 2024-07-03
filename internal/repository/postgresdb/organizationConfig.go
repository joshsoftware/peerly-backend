package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	ae "github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

const (
	createOrganizationQuery = `INSERT INTO organization_config (
		id,
		reward_multiplier,
		reward_quota_renewal_frequency,
		timezone,
		created_by,updated_by)
		VALUES ($1, $2, $3, $4,$5,$6) RETURNING id`

	getOrganizationQuery = `SELECT id,
		reward_multiplier,
		reward_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at,updated_by FROM organization_config WHERE id=$1`

	getOrganizationByIDQuery = `SELECT id,
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

func (s *OrganizationStore) CreateOrganizationConfig(ctx context.Context, org dto.OrganizationConfig) (createdOrganization repository.OrganizationConfig, err error) {
	lastInsertID := 0
	err = s.DB.QueryRow(
		createOrganizationQuery,
		1,
		org.RewardMultiplier,
		org.RewardQuotaRenewalFrequency,
		org.Timezone,
		org.CreatedBy,
		org.UpdatedBy,
	).Scan(&lastInsertID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error creating organization")
		return
	}

	err = s.DB.Get(&createdOrganization, getOrganizationQuery, lastInsertID)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: Log that we can't find the organization even though it's just been created
			log.Error(ae.ErrRecordNotFound, "Just created an Organization, but can't find it!", err)
		}
	}
	return
}

func (s *OrganizationStore) UpdateOrganizationCofig(ctx context.Context, reqOrganization dto.OrganizationConfig) (updatedOrganization repository.OrganizationConfig, err error) {

	updateFields := []string{}
	args := []interface{}{}
	argID := 1

	if reqOrganization.RewardMultiplier != 0 {
		updateFields = append(updateFields, fmt.Sprintf("reward_multiplier = $%d", argID))
		args = append(args, reqOrganization.RewardMultiplier)
		argID++
	}
	if reqOrganization.RewardQuotaRenewalFrequency != 0 {
		updateFields = append(updateFields, fmt.Sprintf("reward_quota_renewal_frequency = $%d", argID))
		args = append(args, reqOrganization.RewardQuotaRenewalFrequency)
		argID++
	}
	if reqOrganization.Timezone != "" {
		updateFields = append(updateFields, fmt.Sprintf("timezone = $%d", argID))
		args = append(args, reqOrganization.Timezone)
		argID++
	}

	if len(updateFields) > 0 {

		updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argID))
		args = append(args, time.Now().UnixMilli())
		argID++

		updateFields = append(updateFields, fmt.Sprintf("updated_by = $%d", argID))
		args = append(args, reqOrganization.UpdatedBy)
		argID++
		// Append the organization ID for the WHERE clause

		args = append(args, 1)
		updateQuery := fmt.Sprintf("UPDATE organization_config SET %s WHERE id = $%d", strings.Join(updateFields, ", "), argID)
		stmt, err := s.DB.Prepare(updateQuery)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error preparing update statement")
			return repository.OrganizationConfig{}, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(args...)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error executing update statement")
			return repository.OrganizationConfig{}, err
		}
	}

	err = s.DB.Get(&updatedOrganization, getOrganizationQuery, reqOrganization.ID)
	if err != nil {
		log.Error(ae.ErrRecordNotFound, "Cannot find organization id "+fmt.Sprint(reqOrganization.ID), err)
		return
	}

	return
}

// GetOrganization - returns an organization from the database if it exists based on its ID primary key
func (s *OrganizationStore) GetOrganizationConfig(ctx context.Context) (organization repository.OrganizationConfig, err error) {
	err = s.DB.Get(&organization, getOrganizationQuery, 1)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("organizationID", 1).Warn("Organization not found")
			return repository.OrganizationConfig{}, ae.OrganizationNotFound
		}
		logger.WithField("err", err.Error()).Error("Error fetching organization")
		return repository.OrganizationConfig{}, err
	}

	return
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
