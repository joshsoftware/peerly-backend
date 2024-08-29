package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type reportAppreciationStore struct {
	DB *sqlx.DB
}

func NewReportRepo(db *sqlx.DB) repository.ReportAppreciationStorer {
	return &reportAppreciationStore{
		DB: db,
	}
}

const (
	createResolution     = `INSERT INTO resolutions (appreciation_id, reporting_comment, reported_by) VALUES ($1,$2,$3) RETURNING id, appreciation_id, reporting_comment, reported_by, reported_at`
	getSenderAndReceiver = `SELECT sender, receiver FROM appreciations WHERE id = $1`
	getReportsCount      = `SELECT count(*) FROM resolutions WHERE appreciation_id = $1 AND reported_by = $2;`
	getAppreciationById  = `SELECT count(*) FROM appreciations WHERE id = $1`
)

func (rs *reportAppreciationStore) CheckAppreciation(ctx context.Context, reqData dto.ReportAppreciationReq) (doesExist bool, err error) {

	var count int64
	doesExist = true
	err = rs.DB.Get(
		&count,
		getAppreciationById,
		reqData.AppreciationId,
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in retriving appreciation count")
		return
	}
	if count == 0 {
		doesExist = false
	}

	return
}

func (rs *reportAppreciationStore) CheckDuplicateReport(ctx context.Context, reqData dto.ReportAppreciationReq) (isDupliate bool, err error) {

	var count int64
	isDupliate = false
	err = rs.DB.Get(
		&count,
		getReportsCount,
		reqData.AppreciationId,
		reqData.ReportedBy,
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in looking for duplicate report")
		return
	}
	if count > 0 {
		isDupliate = true
	}

	return
}

func (rs *reportAppreciationStore) GetSenderAndReceiver(ctx context.Context, reqData dto.ReportAppreciationReq) (resp dto.GetSenderAndReceiverResp, err error) {

	err = rs.DB.GetContext(
		ctx,
		&resp,
		getSenderAndReceiver,
		reqData.AppreciationId,
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in fetching appreciation sender and receiver")
		return
	}
	return
}

func (rs *reportAppreciationStore) ReportAppreciation(ctx context.Context, reportReq dto.ReportAppreciationReq) (resp dto.ReportAppricaitionResp, err error) {

	err = rs.DB.GetContext(
		ctx,
		&resp,
		createResolution,
		reportReq.AppreciationId,
		reportReq.ReportingComment,
		reportReq.ReportedBy,
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in creating report")
		return
	}
	return
}

func (rs *reportAppreciationStore) ListReportedAppreciations(ctx context.Context) (reportedAppreciations []repository.ListReportedAppreciations, err error) {
	query := `select resolutions.id, appreciations.id as appreciation_id, core_values.name as core_value_name, core_values.description as core_value_description, appreciations.description as appreciation_description, appreciations.total_reward_points, appreciations.quarter, appreciations.sender, appreciations.receiver, appreciations.created_at, appreciations.is_valid, resolutions.reporting_comment, resolutions.reported_by, resolutions.reported_at, resolutions.moderator_comment, resolutions.moderated_by, resolutions.moderated_at, resolutions.status from resolutions join appreciations on resolutions.appreciation_id = appreciations.id join core_values on appreciations.core_value_id = core_values.id group by resolutions.id, appreciations.id, core_values.id`
	err = rs.DB.SelectContext(
		ctx,
		&reportedAppreciations,
		query,
	)
	if err != nil {
		err = fmt.Errorf("error in retriving reported appriciations, err:%w", err)
		return
	}
	return
}

func (rs *reportAppreciationStore) CheckResolution(ctx context.Context, id int64) (doesExist bool, appreciation_id int64, err error) {
	query := `select appreciation_id from resolutions where id = $1`
	err = rs.DB.GetContext(
		ctx,
		&appreciation_id,
		query,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			doesExist = false
			err = nil
			return
		}
		err = fmt.Errorf("error in select resolution query, id:%d, err:%w", id, err)
		return
	}
	doesExist = true
	return
}

func (rs *reportAppreciationStore) DeleteAppreciation(ctx context.Context, moderationReq dto.ModerationReq) (err error) {
	moderationQuery := `update resolutions set moderator_comment = $1, moderated_by = $2, status = 'deleted' where id = $3`
	_, err = rs.DB.ExecContext(
		ctx,
		moderationQuery,
		moderationReq.ModeratorComment,
		moderationReq.ModeratedBy,
		moderationReq.ResolutionId,
	)
	if err != nil {
		err = fmt.Errorf("error in updating moderation values, err: %w", err)
		return
	}
	deleteAppreciation := `update appreciations set is_valid = false where id = $1`
	_, err = rs.DB.ExecContext(
		ctx,
		deleteAppreciation,
		moderationReq.AppreciationId,
	)
	if err != nil {
		err = fmt.Errorf("error in marking appreciation invalid, err: %w", err)
		return
	}
	return
}

func (rs *reportAppreciationStore) ResolveAppreciation(ctx context.Context, moderationReq dto.ModerationReq) (err error) {
	moderationQuery := `update resolutions set moderator_comment = $1, moderated_by = $2, status = 'resolved' where id = $3`
	_, err = rs.DB.ExecContext(
		ctx,
		moderationQuery,
		moderationReq.ModeratorComment,
		moderationReq.ModeratedBy,
		moderationReq.ResolutionId,
	)
	if err != nil {
		err = fmt.Errorf("error in updating moderation values, err: %w", err)
		return
	}

	return
}

func (rs *reportAppreciationStore) GetResolution(ctx context.Context, id int64) (reportedAppreciation repository.ListReportedAppreciations, err error) {
	query := `select resolutions.id, appreciations.id as appreciation_id, appreciations.description as appreciation_description, appreciations.total_reward_points, appreciations.quarter, appreciations.sender, appreciations.receiver, appreciations.created_at, appreciations.is_valid, resolutions.reporting_comment, resolutions.reported_by, resolutions.reported_at, resolutions.moderator_comment, resolutions.moderated_by, resolutions.moderated_at, resolutions.status from resolutions join appreciations on resolutions.appreciation_id = appreciations.id where resolutions.id = $1`
	err = rs.DB.GetContext(
		ctx,
		&reportedAppreciation,
		query,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf("no such resolution exists")
			err = apperrors.InvalidId
			return
		}
		logger.Errorf("error in retriving reported appriciation, err:%s", err.Error())
		err = apperrors.InternalServerError
		return
	}
	return
}
