package repository

import (
	"context"
	"database/sql"
	"fmt"

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
		INSERT INTO borrowing_records (book_id, user_id, due_date)
		VALUES ($1, $2, $3)`,
		record.Book.ID, record.UserID, record.DueDate,
	)
	return err
}

func (r *BorrowingRecordRepository) GetBorrowingRecordByID(ctx context.Context, id uuid.UUID) (*models.BorrowingRecord, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, book_id, user_id, borrowed_at, due_date, returned_at FROM borrowing_records WHERE id = $1`, id)
	var record models.BorrowingRecord
	err := row.Scan(&record.ID, &record.Book.ID, &record.UserID, &record.BorrowedAt, &record.DueDate, &record.ReturnedAt)
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
		record.Book.ID, record.UserID, record.BorrowedAt, record.DueDate, record.ReturnedAt, record.ID,
	)
	return err
}

func (r *BorrowingRecordRepository) DeleteBorrowingRecord(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM borrowing_records WHERE id = $1`, id)
	return err
}

func (r *BorrowingRecordRepository) ListBorrowingRecords(ctx context.Context, userID uuid.UUID, queries map[string]string) ([]models.BorrowingRecord, error) {
	baseQuery := `
        SELECT 
            br.id, br.user_id, br.borrowed_at, br.due_date, br.returned_at,
            b.id, b.title, b.author, b.isbn, b.published_date, b.category_id, b.stock, b.added_by, b.created_at, b.updated_at, b.version
        FROM 
            borrowing_records br
        INNER JOIN 
            books b ON br.book_id = b.id
        WHERE 
            br.user_id = $1
    `
	args := []interface{}{userID}
	argPosition := 2

	// Tambahkan kondisi pencarian jika ada
	if queries["title"] != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(b.title) LIKE LOWER($%d)", argPosition)
		args = append(args, "%"+queries["title"]+"%")
		argPosition++
	}

	if queries["status"] == "returned" {
		baseQuery += " AND br.returned_at is not null"
	} else if queries["status"] == "borrowed" {
		baseQuery += "AND br.returned_at is null"
	}

	if queries["order"] != "" {
		baseQuery += " ORDER BY br.borrowed_at " + queries["order"]
	} else {
		baseQuery += " ORDER BY br.borrowed_at desc"
	}

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.BorrowingRecord

	for rows.Next() {
		var record models.BorrowingRecord
		var book models.Book
		var returnedAt sql.NullTime

		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.BorrowedAt,
			&record.DueDate,
			&returnedAt,
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.PublishedDate,
			&book.CategoryID,
			&book.Stock,
			&book.AddedBy,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.Version,
		)
		if err != nil {
			return nil, err
		}

		record.Book = book

		if returnedAt.Valid {
			record.ReturnedAt = &returnedAt.Time
		} else {
			record.ReturnedAt = nil
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
