package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

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

func (r *bookCategoryRepository) Create(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	query := `INSERT INTO book_categories (name) VALUES ($1) returning id`

	if err := r.db.QueryRowContext(ctx, query, name).Scan(&id); err != nil {
		log.Printf("[Repository - Create] Error creating book category: %v", err)
		return uuid.UUID{}, fmt.Errorf("failed to create book category: %w", err)
	}

	return id, nil
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
		log.Printf("[Repository - GetByID] Error getting book category by ID: %v", err)
		return nil, fmt.Errorf("failed to get book category by ID: %w", err)
	}

	return &category, nil
}

func (r *bookCategoryRepository) GetByName(ctx context.Context, name string) (*models.BookCategory, error) {
	query := `SELECT id, name FROM book_categories WHERE name = $1`

	row := r.db.QueryRowContext(ctx, query, name)

	var category models.BookCategory
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("[Repository - GetByName] Error getting book category by Nam: %v", err)
		return nil, fmt.Errorf("failed to get book category by Name: %w", err)
	}

	return &category, nil
}

func (r *bookCategoryRepository) GetAll(ctx context.Context) ([]models.BookCategory, error) {
	query := `SELECT id, name FROM book_categories`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[Repository - GetAll] Error getting all book categories: %v", err)
		return nil, fmt.Errorf("failed to get all book categories: %w", err)
	}
	defer rows.Close()

	var categories []models.BookCategory
	for rows.Next() {
		var category models.BookCategory
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Printf("[Repository - GetAll] Error scanning book category: %v", err)
			return nil, fmt.Errorf("failed to scan book category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[Repository - GetAll] Error during rows iteration: %v", err)
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return categories, nil
}

func (r *bookCategoryRepository) Update(ctx context.Context, category *models.BookCategory) error {
	query := `UPDATE book_categories SET name = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, category.Name, category.ID)
	if err != nil {
		log.Printf("[Repository - Update] Error updating book category: %v", err)
		return fmt.Errorf("failed to update book category: %w", err)
	}

	return nil
}

func (r *bookCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM book_categories WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("[Repository - Delete] Error deleting book category: %v", err)
		return fmt.Errorf("failed to delete book category: %w", err)
	}

	return nil
}
