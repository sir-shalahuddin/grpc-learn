package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
)

type BorrowingRecordRepository struct {
	db *sql.DB
}

func NewBorrowingRecordRepository(db *sql.DB) *BorrowingRecordRepository {
	return &BorrowingRecordRepository{db: db}
}

func (r *BorrowingRecordRepository) CreateBorrowingRecord(ctx context.Context, tx *sql.Tx, record *models.BorrowingRecord) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO borrowing_records (id, book_id, user_id, borrowed_at, due_date, returned_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		record.ID, record.BookID, record.UserID, record.BorrowedAt, record.DueDate, record.ReturnedAt,
	)
	return err
}

func (r *BorrowingRecordRepository) GetBorrowingRecordByID(ctx context.Context, id uuid.UUID) (*models.BorrowingRecord, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, book_id, user_id, borrowed_at, due_date, returned_at FROM borrowing_records WHERE id = $1`, id)
	var record models.BorrowingRecord
	err := row.Scan(&record.ID, &record.BookID, &record.UserID, &record.BorrowedAt, &record.DueDate, &record.ReturnedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *BorrowingRecordRepository) UpdateBorrowingRecord(ctx context.Context, tx *sql.Tx, record *models.BorrowingRecord) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE borrowing_records 
		SET book_id = $1, user_id = $2, borrowed_at = $3, due_date = $4, returned_at = $5
		WHERE id = $6`,
		record.BookID, record.UserID, record.BorrowedAt, record.DueDate, record.ReturnedAt, record.ID,
	)
	return err
}

func (r *BorrowingRecordRepository) DeleteBorrowingRecord(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM borrowing_records WHERE id = $1`, id)
	return err
}

func (r *BorrowingRecordRepository) ListBorrowingRecordsByBookID(ctx context.Context, bookID uuid.UUID) ([]models.BorrowingRecord, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, book_id, user_id, borrowed_at, due_date, returned_at FROM borrowing_records WHERE book_id = $1`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.BorrowingRecord
	for rows.Next() {
		var record models.BorrowingRecord
		if err := rows.Scan(&record.ID, &record.BookID, &record.UserID, &record.BorrowedAt, &record.DueDate, &record.ReturnedAt); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
