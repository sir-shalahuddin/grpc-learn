package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/models"
)

type bookCategoryRepository struct {
	db *sql.DB
}

// NewBookCategoryRepository returns a new instance of a BookCategoryRepository.
func NewBookCategoryRepository(db *sql.DB) *bookCategoryRepository {
	return &bookCategoryRepository{
		db: db,
	}
}

func (r *bookCategoryRepository) Create(ctx context.Context, category *models.BookCategory) error {
	query := `INSERT INTO book_categories (id, name) VALUES ($1, $2)`

	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, query, category.ID, category.Name)
	if err != nil {
		return fmt.Errorf("failed to create book category: %w", err)
	}

	return nil
}

func (r *bookCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.BookCategory, error) {
	query := `SELECT id, name FROM book_categories WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var category models.BookCategory
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get book category by ID: %w", err)
	}

	return &category, nil
}

func (r *bookCategoryRepository) GetAll(ctx context.Context) ([]models.BookCategory, error) {
	query := `SELECT id, name FROM book_categories`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all book categories: %w", err)
	}
	defer rows.Close()

	var categories []models.BookCategory
	for rows.Next() {
		var category models.BookCategory
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan book category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return categories, nil
}

func (r *bookCategoryRepository) Update(ctx context.Context, category *models.BookCategory) error {
	query := `UPDATE book_categories SET name = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, category.Name, category.ID)
	if err != nil {
		return fmt.Errorf("failed to update book category: %w", err)
	}

	return nil
}

func (r *bookCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM book_categories WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book category: %w", err)
	}

	return nil
}
