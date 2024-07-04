package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
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

// insert into appreciations
// (core_value_id, description, is_valid, quarter, sender, receiver)
// values
// (2, 'desc', true, 1, 1153, 1154);
