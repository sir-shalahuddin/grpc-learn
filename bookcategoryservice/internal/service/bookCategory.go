package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/models"
)

type bookCategoryRepository interface {
	Create(ctx context.Context, category *models.BookCategory) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.BookCategory, error)
	GetAll(ctx context.Context) ([]models.BookCategory, error)
	Update(ctx context.Context, category *models.BookCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type bookCategoryService struct {
	repo bookCategoryRepository
}

// NewBookCategoryService returns a new instance of BookCategoryService.
func NewBookCategoryService(repo bookCategoryRepository) *bookCategoryService {
	return &bookCategoryService{
		repo: repo,
	}
}

func (s *bookCategoryService) CreateCategory(ctx context.Context, name string) (uuid.UUID, error) {
	if name == "" {
		return uuid.Nil, errors.New("category name cannot be empty")
	}

	category := &models.BookCategory{
		Name: name,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category.ID, nil
}

func (s *bookCategoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.BookCategory, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid category ID")
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}

	return category, nil
}

func (s *bookCategoryService) GetAllCategories(ctx context.Context) ([]models.BookCategory, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}

	return categories, nil
}

func (s *bookCategoryService) UpdateCategory(ctx context.Context, category *models.BookCategory) error {
	if category.ID == uuid.Nil {
		return errors.New("invalid category ID")
	}
	if category.Name == "" {
		return errors.New("category name cannot be empty")
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

func (s *bookCategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid category ID")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
