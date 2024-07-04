package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type BaseRepository struct {
	DB *sqlx.DB
}

type BaseTransaction struct {
	tx *sqlx.Tx
}

func (repo *BaseRepository) BeginTx(ctx context.Context) (repository.Transaction, error) {

	txObj, err := repo.DB.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("error occured while initiating database transaction: %v", err.Error())
		return nil, err
	}

	return &BaseTransaction{
		tx: txObj,
	}, nil
}

func (repo *BaseRepository) HandleTransaction(ctx context.Context, tx repository.Transaction, isSuccess bool) error {    var err error
    if !isSuccess {
        err = tx.Rollback()
        if err != nil {
            log.Printf("error occurred while rollback database transaction: %v", err.Error())
            return err
        }
        return err
    }

    err = tx.Commit()
    if err != nil {
        log.Printf("error occurred while commit database transaction: %v", err.Error())
        return err
    }
    return err
}

func (repo *BaseTransaction) Commit() error {
	return repo.tx.Commit()
}

func (repo *BaseTransaction) Rollback() error {
	return repo.tx.Rollback()
}

func (repo *BaseRepository) InitiateQueryExecutor(tx repository.Transaction) (executor sqlx.Ext) {
	//Populate the query executor so we can use a transaction if one is present.
	//If we are not running inside a transaction then the plain sqlx.DB object is used.
	executor = repo.DB
	if tx != nil {
		txObj := tx.(*BaseTransaction)
		executor = txObj.tx
	}

	return executor
}
