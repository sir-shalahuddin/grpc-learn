package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) CreateBook(ctx context.Context, book *models.Book) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO books (id, title, author, isbn, published_date, category_id, stock, added_by, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		book.ID, book.Title, book.Author, book.ISBN, book.PublishedDate, book.CategoryID, book.Stock, book.AddedBy, book.CreatedAt, book.UpdatedAt, book.Version,
	)
	return err
}

func (r *BookRepository) GetBookByID(ctx context.Context, id uuid.UUID) (*models.Book, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, title, author, isbn, published_date, category_id, stock, added_by, created_at, updated_at, version FROM books WHERE id = $1`, id)
	var book models.Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedDate, &book.CategoryID, &book.Stock, &book.AddedBy, &book.CreatedAt, &book.UpdatedAt, &book.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get book by ID: %w", err)
	}
	return &book, nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, isbn = $3, published_date = $4, category_id = $5, stock = $6, added_by = $7, updated_at = $8, version = version + 1
		WHERE id = $9 AND version = $10
	`
	var result sql.Result
	var err error

	if tx != nil {
		// If a transaction is provided, use it.
		result, err = tx.ExecContext(ctx, query, book.Title, book.Author, book.ISBN, book.PublishedDate, book.CategoryID, book.Stock, book.AddedBy, time.Now(), book.ID, book.Version)
	} else {
		// Otherwise, use the repository's database connection.
		result, err = r.db.ExecContext(ctx, query, book.Title, book.Author, book.ISBN, book.PublishedDate, book.CategoryID, book.Stock, book.AddedBy, time.Now(), book.ID, book.Version)
	}

	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conflict detected: book record was modified by another user")
	}

	return nil
}

func (r *BookRepository) DeleteBook(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}

func (r *BookRepository) ListBooks(ctx context.Context, title, author, category string) ([]*models.Book, error) {
	query := `SELECT id, title, author, isbn, published_date, category_id, stock, added_by, created_at, updated_at, version FROM books WHERE 1=1`
	args := []interface{}{}
	if title != "" {
		query += ` AND title LIKE $1`
		args = append(args, "%"+title+"%")
	}
	if author != "" {
		query += ` AND author LIKE $2`
		args = append(args, "%"+author+"%")
	}
	if category != "" {
		query += ` AND category_id = $3`
		args = append(args, category)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list books: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedDate, &book.CategoryID, &book.Stock, &book.AddedBy, &book.CreatedAt, &book.UpdatedAt, &book.Version); err != nil {
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, &book)
	}
	return books, nil
}
