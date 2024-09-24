package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/dto"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/models"
)

type bookCategoryRepository interface {
	Create(ctx context.Context, name string) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.BookCategory, error)
	GetByName(ctx context.Context, name string) (*models.BookCategory, error)
	GetAll(ctx context.Context) ([]models.BookCategory, error)
	Update(ctx context.Context, category *models.BookCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type bookCategoryService struct {
	repo bookCategoryRepository
}

var (
	ErrDuplicateCategory = errors.New("category already exist")
)

// NewBookCategoryService returns a new instance of BookCategoryService.
func NewBookCategoryService(repo bookCategoryRepository) *bookCategoryService {
	return &bookCategoryService{
		repo: repo,
	}
}

func (s *bookCategoryService) CreateCategory(ctx context.Context, req dto.CreateBookCategoryRequest) (dto.CreateBookCategoryResponse, error) {
	existCategory, err := s.repo.GetByName(ctx, req.Name)
	if err != nil {
		return dto.CreateBookCategoryResponse{}, err
	}
	if existCategory != nil {
		return dto.CreateBookCategoryResponse{}, ErrDuplicateCategory
	}

	res, err := s.repo.Create(ctx, req.Name)
	if err != nil {
		return dto.CreateBookCategoryResponse{}, err
	}
	return dto.CreateBookCategoryResponse{ID: res, Name: req.Name}, nil
}

func (s *bookCategoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.BookCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *bookCategoryService) GetAllCategories(ctx context.Context) ([]models.BookCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *bookCategoryService) UpdateCategory(ctx context.Context, category *models.BookCategory) error {
	return s.repo.Update(ctx, category)
}

func (s *bookCategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
