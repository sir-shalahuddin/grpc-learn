package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
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

type bookService struct {
	bookRepo BookRepository
}

func NewBookService(bookRepo BookRepository) *bookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}

func (s *bookService) CreateBook(ctx context.Context, book *models.Book) error {
	return s.bookRepo.CreateBook(ctx, book)
}

func (s *bookService) GetBookByID(ctx context.Context, id uuid.UUID) (*models.Book, error) {
	return s.bookRepo.GetBookByID(ctx, id)
}

func (s *bookService) UpdateBook(ctx context.Context, book *models.Book) error {
	return s.bookRepo.UpdateBook(ctx, nil, book)
}

func (s *bookService) DeleteBook(ctx context.Context, id uuid.UUID) error {
	return s.bookRepo.DeleteBook(ctx, id)
}

func (s *bookService) ListBooks(ctx context.Context, title, author, category string) ([]*models.Book, error) {
	return s.bookRepo.ListBooks(ctx, title, author, category)
}
