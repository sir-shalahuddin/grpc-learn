package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
)

var (
	ErrBorrowingRecordNotFound = errors.New("borrowing record not found")
	ErrBookUnavailable         = errors.New("book is currently unavailable")
)

type BorrowingRecordRepository interface {
	CreateBorrowingRecord(ctx context.Context, tx *sql.Tx, record *models.BorrowingRecord) error
	GetBorrowingRecordByID(ctx context.Context, id uuid.UUID) (*models.BorrowingRecord, error)
	UpdateBorrowingRecord(ctx context.Context, record *models.BorrowingRecord) error
	DeleteBorrowingRecord(ctx context.Context, id uuid.UUID) error
	ListBorrowingRecordsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.BorrowingRecord, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	GetBookByID(ctx context.Context, bookID uuid.UUID) (*models.Book, error)
	UpdateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error
}

type TxRepository interface {
}

type borrowingRecordService struct {
	repo BorrowingRecordRepository
}

func NewBorrowingRecordService(repo BorrowingRecordRepository) *borrowingRecordService {
	return &borrowingRecordService{
		repo: repo,
	}
}

func (s *borrowingRecordService) BorrowBook(ctx context.Context, bookID, userID uuid.UUID, dueDate *time.Time) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.repo.Rollback(tx)

	book, err := s.repo.GetBookByID(ctx, bookID)
	if err != nil {
		return err
	}
	if book == nil || book.Stock <= 0 {
		return fmt.Errorf("book not available")
	}

	// Create borrowing record
	record := &models.BorrowingRecord{
		ID:         uuid.New(),
		BookID:     bookID,
		UserID:     userID,
		BorrowedAt: time.Now(),
		DueDate:    dueDate,
	}
	if err := s.repo.CreateBorrowingRecord(ctx, tx, record); err != nil {
		return err
	}

	// Update book stock
	book.Stock--
	if err := s.repo.UpdateBook(ctx, tx, book); err != nil {
		return err
	}

	return s.repo.Commit(tx)
}

func (s *borrowingRecordService) ReturnBook(ctx context.Context, recordID uuid.UUID) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.repo.Rollback(tx)

	record, err := s.repo.GetBorrowingRecordByID(ctx, recordID)
	if err != nil {
		return err
	}
	if record == nil {
		return fmt.Errorf("borrowing record not found")
	}

	book, err := s.repo.GetBookByID(ctx, record.BookID)
	if err != nil {
		return err
	}

	// Update record
	record.ReturnedAt = time.Now()
	if err := s.repo.UpdateBorrowingRecord(ctx, record); err != nil {
		return err
	}

	// Update book stock
	book.Stock++
	if err := s.repo.UpdateBook(ctx, nil, book); err != nil {
		return err
	}

	return s.repo.Commit(tx)
}

func (s *borrowingRecordService) ListBorrowingRecordsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.BorrowingRecord, error) {
	return s.repo.ListBorrowingRecordsByBookID(ctx, bookID)
}
