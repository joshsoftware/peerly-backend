package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// RepositoryTransaction interface holds transaction specific methods
type RepositoryTransaction interface {
	// return a transaction from a sql connection
	BeginTx(ctx context.Context) (Transaction, error)
	HandleTransaction(ctx context.Context, tx Transaction, isSuccess bool) error
	InitiateQueryExecutor(tx Transaction) (executor sqlx.Ext)
}

type Transaction interface {
	Commit() error
	Rollback() error
}
