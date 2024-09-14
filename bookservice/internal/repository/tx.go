package repository

import (
	"context"
	"database/sql"
)

type txRepository struct {
	db *sql.DB
}

func NewTxRepository(db *sql.DB) *txRepository {
	return &txRepository{db: db}
}

func (repo *txRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (repo *txRepository) Commit(tx *sql.Tx) error {
	if tx == nil {
		return nil
	}
	return tx.Commit()
}

func (repo *txRepository) Rollback(tx *sql.Tx) error {
	if tx == nil {
		return nil
	}
	return tx.Rollback()
}
