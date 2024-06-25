package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	ae "github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
	"github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

const (
	createOrganizationQuery = `INSERT INTO organizations (
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	deleteOrganizationQuery = `UPDATE organizations SET soft_delete = true, soft_delete_by = $1 WHERE id = $2`

	getOrganizationQuery = `SELECT id,
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at FROM organizations WHERE id=$1 AND soft_delete = FALSE`

	listOrganizationsQuery = `SELECT id,
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at FROM organizations WHERE soft_delete = FALSE ORDER BY name ASC`

	getOrganizationByDomainNameQuery = `SELECT id,
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at FROM organizations WHERE domain_name=$1 AND soft_delete = FALSE LIMIT 1`
	getOrganizationByIDQuery = `SELECT id,
		name,
		contact_email,
		domain_name,
		subscription_status,
		subscription_valid_upto,
		hi5_limit,
		hi5_quota_renewal_frequency,
		timezone,
		created_at,
		created_by,
		updated_at FROM organizations WHERE id=$1 LIMIT 1`
	getCountOfContactEmailQuery = `SELECT COUNT(*) FROM organizations WHERE contact_email = $1 AND soft_delete = FALSE`
	getCountOfDomainNameQuery   = `SELECT COUNT(*) FROM organizations WHERE domain_name = $1 AND soft_delete = FALSE`
	getCountOfIdQuery           = `SELECT COUNT(*) FROM organizations WHERE id = $1 AND soft_delete = FALSE`
)

type OrganizationStore struct {
	BaseRepository
}

func NewOrganizationRepo(db *sqlx.DB) repository.OrganizationStorer {
	return &OrganizationStore{
		BaseRepository: BaseRepository{db}, // Use *sqlx.DB instead of *sql.DB
	}
}

func (orgStr *OrganizationStore) ListOrganizations(ctx context.Context) (organizations []repository.Organization, err error) {
	err = orgStr.DB.Select(&organizations, listOrganizationsQuery)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error listing organizations")
		return organizations, ae.InternalServer
	}
	return
}

func (s *OrganizationStore) CreateOrganization(ctx context.Context, org dto.Organization) (createdOrganization repository.Organization, err error) {
	// Set org.CreatedAt so we get a valid created_at value from the database going forward
	org.CreatedAt = time.Now().UTC()

	lastInsertID := 0
	err = s.DB.QueryRow(
		createOrganizationQuery,
		org.Name,
		org.ContactEmail,
		org.DomainName,
		org.SubscriptionStatus,
		org.SubscriptionValidUpto,
		org.Hi5Limit,
		org.Hi5QuotaRenewalFrequency,
		org.Timezone,
		org.CreatedBy,
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

func (s *OrganizationStore) UpdateOrganization(ctx context.Context, reqOrganization dto.Organization) (updatedOrganization repository.Organization, err error) {
	err = s.DB.Get(&updatedOrganization, getOrganizationQuery, reqOrganization.ID)
	if err != nil {
		log.Error(ae.ErrRecordNotFound, "Cannot find organization id "+fmt.Sprint(reqOrganization.ID), err)
		return repository.Organization{}, ae.OrganizationNotFound
	}

	var dbOrganization repository.Organization
	err = s.DB.Get(&dbOrganization, getOrganizationQuery, reqOrganization.ID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching organization")
		return
	}

	updateFields := []string{}
	args := []interface{}{}
	argID := 1

	if reqOrganization.Name != "" {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", argID))
		args = append(args, reqOrganization.Name)
		argID++
	}
	if reqOrganization.ContactEmail != "" {
		updateFields = append(updateFields, fmt.Sprintf("contact_email = $%d", argID))
		updateFields = append(updateFields,"is_email_verified = false")
		args = append(args, reqOrganization.ContactEmail)
		argID++
	}
	if reqOrganization.DomainName != "" {
		updateFields = append(updateFields, fmt.Sprintf("domain_name = $%d", argID))
		args = append(args, reqOrganization.DomainName)
		argID++
	}

	if !reqOrganization.SubscriptionValidUpto.IsZero() {
		updateFields = append(updateFields, fmt.Sprintf("subscription_valid_upto = $%d", argID))
		args = append(args, reqOrganization.SubscriptionValidUpto)
		argID++
		updateFields = append(updateFields, fmt.Sprintf("subscription_status = $%d", argID))
		args = append(args, 1)
		argID++
	}
	if reqOrganization.Hi5Limit != 0 {
		updateFields = append(updateFields, fmt.Sprintf("hi5_limit = $%d", argID))
		args = append(args, reqOrganization.Hi5Limit)
		argID++
	}
	if reqOrganization.Hi5QuotaRenewalFrequency != "" {
		updateFields = append(updateFields, fmt.Sprintf("hi5_quota_renewal_frequency = $%d", argID))
		args = append(args, reqOrganization.Hi5QuotaRenewalFrequency)
		argID++
	}
	if reqOrganization.Timezone != "" {
		updateFields = append(updateFields, fmt.Sprintf("timezone = $%d", argID))
		args = append(args, reqOrganization.Timezone)
		argID++
	}

	if len(updateFields) > 0 {

		updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argID))
		args = append(args, time.Now())
		argID++
		// Append the organization ID for the WHERE clause

		args = append(args, reqOrganization.ID)
		updateQuery := fmt.Sprintf("UPDATE organizations SET %s WHERE id = $%d", strings.Join(updateFields, ", "), argID)
		fmt.Println("update query: ------------->\n", updateQuery)
		fmt.Println("update args: ------------->\n", args)
		stmt, err := s.DB.Prepare(updateQuery)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error preparing update statement")
			return repository.Organization{}, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(args...)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error executing update statement")
			return repository.Organization{}, err
		}
	}

	err = s.DB.Get(&updatedOrganization, getOrganizationQuery, reqOrganization.ID)
	if err != nil {
		log.Error(ae.ErrRecordNotFound, "Cannot find organization id "+fmt.Sprint(reqOrganization.ID), err)
		return
	}

	return
}

func (s *OrganizationStore) DeleteOrganization(ctx context.Context, organizationID int, userId int64) (err error) {
	sqlRes, err := s.DB.Exec(deleteOrganizationQuery, userId, organizationID)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting organization")
		return ae.InternalServer
	}

	rowsAffected, err := sqlRes.RowsAffected()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching rows affected count")
		return ae.InternalServer
	}

	if rowsAffected == 0 {
		err = fmt.Errorf("organization with ID %d not found", organizationID)
		logger.WithField("organizationID", organizationID).Warn(err.Error())
		return ae.OrganizationNotFound
	}

	return nil
}

// GetOrganization - returns an organization from the database if it exists based on its ID primary key
func (s *OrganizationStore) GetOrganization(ctx context.Context, organizationID int) (organization repository.Organization, err error) {
	err = s.DB.Get(&organization, getOrganizationQuery, organizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("organizationID", organizationID).Warn("Organization not found")
			return repository.Organization{}, ae.OrganizationNotFound
		}
		logger.WithField("err", err.Error()).Error("Error fetching organization")
		return repository.Organization{}, err
	}

	return
}

func (s *OrganizationStore) GetOrganizationByDomainName(ctx context.Context, domainName string) (organization repository.Organization, err error) {
	fmt.Println("GetOrganizationByDomainName ------------------------>")
	err = s.DB.Get(&organization, getOrganizationByDomainNameQuery, domainName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithField("organization domain name", domainName).Warn("Organization not found by domain name")
			return repository.Organization{}, ae.OrganizationNotFound
		}
		logger.WithField("err", err.Error()).Error("Error fetching organization")
		return repository.Organization{}, err
	}
	return
}

func (s *OrganizationStore) IsEmailPresent(ctx context.Context, email string) bool {

	var count int

	err := s.DB.QueryRowContext(ctx, getCountOfContactEmailQuery, email).Scan(&count)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching contact email of organization by contact email id: " + email)
		return false
	}

	return count > 0
}

func (s *OrganizationStore) IsDomainPresent(ctx context.Context, domainName string) bool {

	var count int

	err := s.DB.QueryRowContext(ctx, getCountOfDomainNameQuery, domainName).Scan(&count)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching domain name of organization by contact email id: " + domainName)
		return false
	}

	return count > 0
}

func (s *OrganizationStore) IsOrganizationIdPresent(ctx context.Context, organizationId int64) bool {
	var count int

	err := s.DB.QueryRowContext(ctx, getCountOfIdQuery, organizationId).Scan(&count)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching id of organization: " + strconv.FormatInt(organizationId, 10))
		return false
	}

	return count > 0
}

///helper functions Organization

func OrganizationToDB(org dto.Organization) repository.Organization {
	return repository.Organization{
		ID:                       org.ID,
		Name:                     org.Name,
		ContactEmail:             org.ContactEmail,
		DomainName:               org.DomainName,
		SubscriptionStatus:       org.SubscriptionStatus,
		SubscriptionValidUpto:    org.SubscriptionValidUpto,
		Hi5Limit:                 org.Hi5Limit,
		Hi5QuotaRenewalFrequency: org.Hi5QuotaRenewalFrequency,
		Timezone:                 org.Timezone,
		CreatedAt:                org.CreatedAt,
	}
}