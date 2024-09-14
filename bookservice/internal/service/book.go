package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
	pb "github.com/sir-shalahuddin/grpc-learn/bookservice/proto/categoryservice"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type BookRepository interface {
	GetBookByID(ctx context.Context, bookID uuid.UUID) (*models.Book, error)
	CreateBook(ctx context.Context, book *models.Book) error
	UpdateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error
	DeleteBook(ctx context.Context, bookID uuid.UUID) error
	ListBooks(ctx context.Context, title, author, category string) ([]*models.Book, error)
}

type categoryRepository interface {
	GetCategories(ctx context.Context) ([]*pb.CategoryResponse, error)
	GetCategoryByID(ctx context.Context, id string) (*pb.CategoryResponse, error)
}

type bookService struct {
	bookRepo BookRepository
	ctgRepo  categoryRepository
}

func NewBookService(bookRepo BookRepository, ctgRepo categoryRepository) *bookService {
	return &bookService{
		bookRepo: bookRepo,
		ctgRepo:  ctgRepo,
	}
}

func (s *bookService) CreateBook(ctx context.Context, book *models.Book) error {
	return s.bookRepo.CreateBook(ctx, book)
}

func (s *bookService) GetBookByID(ctx context.Context, id uuid.UUID) (*models.Book, error) {
	book, err := s.bookRepo.GetBookByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get book by ID: %w", err)
	}
	if book == nil {
		return nil, nil
	}

	// Fetch the category name using the category ID
	category, err := s.ctgRepo.GetCategoryByID(ctx, book.CategoryID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}
	if category != nil {
		book.CategoryName = category.Name
	}

	return book, nil
}

func (s *bookService) UpdateBook(ctx context.Context, book *models.Book) error {
	return s.bookRepo.UpdateBook(ctx, nil, book)
}

func (s *bookService) DeleteBook(ctx context.Context, id uuid.UUID) error {
	return s.bookRepo.DeleteBook(ctx, id)
}

func (s *bookService) ListBooks(ctx context.Context, title, author, category string) ([]*models.Book, error) {
	// Get books from the local repository
	books, err := s.bookRepo.ListBooks(ctx, title, author, category)
	if err != nil {
		return nil, fmt.Errorf("failed to list books from repository: %w", err)
	}

	// Get all categories from the BookCategoryService using gRPC client
	categoryResp, err := s.ctgRepo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories from gRPC service: %w", err)
	}

	// Create a map to easily find category names by ID
	categoryMap := make(map[string]string)
	for _, c := range categoryResp {
		categoryMap[c.Id] = c.Name
	}

	// Map category names to books
	for _, book := range books {
		if categoryName, ok := categoryMap[book.CategoryID.String()]; ok {
			book.CategoryName = categoryName
		} else {
			book.CategoryName = "Unknown" 
		}
	}

	return books, nil
}
