package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
	pb "github.com/sir-shalahuddin/grpc-learn/bookservice/proto/categoryservice"
)

var (
	ErrBookNotFound     = errors.New("book not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrBookDuplicate    = errors.New("book already exists")
)

type BookRepository interface {
	GetBookByID(ctx context.Context, bookID uuid.UUID) (*models.Book, error)
	AddBook(ctx context.Context, book *models.Book) error
	UpdateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error
	DeleteBook(ctx context.Context, bookID uuid.UUID) error
	ListBooks(ctx context.Context, title, author, category string, limit, offset int) ([]*models.Book, error)
	GetBookByISBN(ctx context.Context, isbn string) (*models.Book, error)
}

type categoryRepository interface {
	GetCategories(ctx context.Context) ([]*pb.CategoryResponse, error)
	GetCategoryByID(ctx context.Context, id string) (*pb.CategoryResponse, error)
}

type bookService struct {
	bookRepo      BookRepository
	ctgRepo       categoryRepository
	categoryCache []*pb.CategoryResponse
	lastUpdate    time.Time
	categoryMu    sync.Mutex
}

func NewBookService(bookRepo BookRepository, ctgRepo categoryRepository) *bookService {
	return &bookService{
		bookRepo: bookRepo,
		ctgRepo:  ctgRepo,
	}
}

func (s *bookService) AddBook(ctx context.Context, req dto.AddBookRequest, userID uuid.UUID) error {

	req.ISBN = strings.ReplaceAll(req.ISBN, "-", "")

	if req.ISBN != "" {
		existingBook, err := s.bookRepo.GetBookByISBN(ctx, req.ISBN)
		if err != nil {
			return err
		}
		if existingBook != nil {
			return fmt.Errorf("%w: %s", ErrBookDuplicate, existingBook.Title)
		}
	}

	if req.CategoryID != uuid.Nil {
		existingCategory, err := s.ctgRepo.GetCategoryByID(ctx, req.CategoryID.String())
		if err != nil {
			return err
		}
		if existingCategory == nil {
			return ErrCategoryNotFound
		}
		// log.Println("hit", existingCategory.GetId())
	}

	book := &models.Book{
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.ISBN,
		PublishedDate: req.PublishedDate,
		CategoryID:    req.CategoryID,
		Stock:         req.Stock,
		AddedBy:       userID,
	}

	return s.bookRepo.AddBook(ctx, book)
}

func (s *bookService) GetBookByID(ctx context.Context, id uuid.UUID) (*dto.GetBookResponse, error) {
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

	response := dto.GetBookResponse{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		ISBN:          book.ISBN,
		PublishedDate: book.PublishedDate,
		Category:      category.Name,
		Stock:         book.Stock,
	}

	return &response, nil
}

func (s *bookService) UpdateBook(ctx context.Context, req dto.UpdateBookRequest, bookID uuid.UUID) error {
	req.ISBN = strings.ReplaceAll(req.ISBN, "-", "")

	if req.ISBN != "" {
		existingBook, err := s.bookRepo.GetBookByISBN(ctx, req.ISBN)
		if err != nil {
			return err
		}
		if existingBook != nil {
			return fmt.Errorf("%w: %s", ErrBookDuplicate, existingBook.Title)
		}
	}

	if req.CategoryID != uuid.Nil {
		existingCategory, err := s.ctgRepo.GetCategoryByID(ctx, req.CategoryID.String())
		if err != nil {
			return err
		}
		if existingCategory == nil {
			return ErrCategoryNotFound
		}
	}

	book := &models.Book{
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.ISBN,
		PublishedDate: req.PublishedDate,
		CategoryID:    req.CategoryID,
		Stock:         req.Stock,
		ID:            bookID,
	}

	return s.bookRepo.UpdateBook(ctx, nil, book)
}

func (s *bookService) DeleteBook(ctx context.Context, bookID uuid.UUID) error {
	return s.bookRepo.DeleteBook(ctx, bookID)
}

// ListBooks lists books and maps categories using cache
func (s *bookService) ListBooks(ctx context.Context, title, author, category string, page string) ([]*dto.GetBookResponse, error) {

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limit := 10
	offset := (pageNum - 1) * limit

	var (
		books       []*models.Book
		categoryMap map[string]string
	)

	// Use WaitGroup to run ListBooks and GetCategories concurrently
	var wg sync.WaitGroup
	wg.Add(2) // We are running 2 operations in parallel

	// Channel to capture errors
	errCh := make(chan error, 2)

	// Goroutine to fetch books
	go func() {
		defer wg.Done()
		var bookErr error
		books, bookErr = s.bookRepo.ListBooks(ctx, title, author, category, limit, offset)
		if bookErr != nil {
			errCh <- fmt.Errorf("failed to list books from repository: %w", bookErr)
		}
	}()

	// Goroutine to fetch categories
	go func() {
		defer wg.Done()

		// Fetch categories from the repository
		categories, catErr := s.ctgRepo.GetCategories(ctx)
		if catErr != nil {
			errCh <- fmt.Errorf("failed to get categories from gRPC service: %w", catErr)
			return
		}

		// Create a map to easily find category names by ID
		categoryMap = make(map[string]string)
		for _, c := range categories {
			categoryMap[c.Id] = c.Name
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Close the error channel to avoid leaking
	close(errCh)

	// Check if there were any errors during parallel execution
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	// Prepare response slice after both goroutines have completed
	var responses []*dto.GetBookResponse
	for _, book := range books {
		categoryName := "Unknown"
		if name, ok := categoryMap[book.CategoryID.String()]; ok {
			categoryName = name
		}

		response := dto.GetBookResponse{
			ID:            book.ID,
			Title:         book.Title,
			Author:        book.Author,
			ISBN:          book.ISBN,
			PublishedDate: book.PublishedDate,
			Category:      categoryName,
			Stock:         book.Stock,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}
