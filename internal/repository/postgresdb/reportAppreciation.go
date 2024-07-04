package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type reportAppreciationStore struct {
	DB *sqlx.DB
}

func NewReportRepo(db *sqlx.DB) repository.ReportAppreciationStorer {
	return &reportAppreciationStore{
		DB: db,
	}
}

func (rs *reportAppreciationStore) ReportAppreciation(ctx context.Context, reportReq dto.ReportAppreciationReq) (err error) {

	return
}
