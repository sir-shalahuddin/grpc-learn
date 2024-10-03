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

func (r *BookRepository) AddBook(ctx context.Context, book *models.Book) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO books (title, author, isbn, published_date, category_id, stock, added_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		book.Title, book.Author, book.ISBN, book.PublishedDate, book.CategoryID, book.Stock, book.AddedBy,
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

func (r *BookRepository) GetBookByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, title, author, isbn, published_date, category_id, stock, added_by, created_at, updated_at, version FROM books WHERE isbn = $1`, isbn)
	var book models.Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedDate, &book.CategoryID, &book.Stock, &book.AddedBy, &book.CreatedAt, &book.UpdatedAt, &book.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get book by ISBN: %w", err)
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
		result, err = tx.ExecContext(ctx, query, book.Title, book.Author, book.ISBN, book.PublishedDate, book.CategoryID, book.Stock, book.AddedBy, time.Now(), book.ID, book.Version)
	} else {
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

func (r *BookRepository) ListBooks(ctx context.Context, title, author, category string, limit, offset int) ([]*models.Book, error) {
	query := `SELECT id, title, author, isbn, published_date, category_id, stock, added_by, created_at, updated_at, version FROM books WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if title != "" {
		query += fmt.Sprintf(" AND title LIKE $%d", argIndex)
		args = append(args, "%"+title+"%")
		argIndex++
	}
	if author != "" {
		query += fmt.Sprintf(" AND author LIKE $%d", argIndex)
		args = append(args, "%"+author+"%")
		argIndex++
	}
	if category != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	// Add pagination: LIMIT and OFFSET
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list books: %w", err)
	}
	defer rows.Close()

	// Parse the result set
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
