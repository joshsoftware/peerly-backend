package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

// BaseRepository holds Database instance
type BaseRepository struct {
	DB *sqlx.DB
}

// BaseTransaction holds transaction instance
type BaseTransaction struct {
	tx *sqlx.Tx
}

// BeginTx begins transaction and return transaction instance
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

// HandleTransaction commit transaction when transaction is successful else rollback
func (repo *BaseRepository) HandleTransaction(ctx context.Context, tx repository.Transaction, isSuccess bool) error {
	var err error
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

// Commit method commit the transaction
func (repo *BaseTransaction) Commit() error {
	return repo.tx.Commit()
}

// Rollback method Rollback the transaction
func (repo *BaseTransaction) Rollback() error {
	return repo.tx.Rollback()
}

// InitiateQueryExecutor Populate the query executor so we can use a transaction if one is present.
// If we are not running inside a transaction then the plain sqlx.DB object is used.
func (repo *BaseRepository) InitiateQueryExecutor(tx repository.Transaction) (executor sqlx.Ext) {

	executor = repo.DB
	if tx != nil {
		txObj := tx.(*BaseTransaction)
		executor = txObj.tx
	}

	return executor
}
