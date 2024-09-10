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
	UpdateBorrowingRecord(ctx context.Context, tx *sql.Tx, record *models.BorrowingRecord) error
	DeleteBorrowingRecord(ctx context.Context, id uuid.UUID) error
	ListBorrowingRecordsByBookID(ctx context.Context, bookID uuid.UUID) ([]models.BorrowingRecord, error)
}

type TxRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
}

type borrowingRecordService struct {
	repo     BorrowingRecordRepository
	txRepo   TxRepository
	bookRepo BookRepository
}

func NewBorrowingRecordService(repo BorrowingRecordRepository, txRepo TxRepository, bookRepo BookRepository) *borrowingRecordService {
	return &borrowingRecordService{
		repo:     repo,
		txRepo:   txRepo,
		bookRepo: bookRepo,
	}
}

func (s *borrowingRecordService) BorrowBook(ctx context.Context, bookID, userID uuid.UUID, dueDate *time.Time) error {
	book, err := s.bookRepo.GetBookByID(ctx, bookID)
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

	tx, err := s.txRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.txRepo.Rollback(tx)

	if err := s.repo.CreateBorrowingRecord(ctx, tx, record); err != nil {
		return err
	}

	// Update book stock
	book.Stock--
	if err := s.bookRepo.UpdateBook(ctx, tx, book); err != nil {
		return err
	}

	return s.txRepo.Commit(tx)
}

func (s *borrowingRecordService) ReturnBook(ctx context.Context, recordID uuid.UUID) error {
	record, err := s.repo.GetBorrowingRecordByID(ctx, recordID)
	if err != nil {
		return err
	}
	if record == nil {
		return ErrBorrowingRecordNotFound
	}

	book, err := s.bookRepo.GetBookByID(ctx, record.BookID)
	if err != nil {
		return err
	}

	// Update record
	record.ReturnedAt = time.Now()
	tx, err := s.txRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.txRepo.Rollback(tx)

	if err := s.repo.UpdateBorrowingRecord(ctx, tx, record); err != nil {
		return err
	}

	// Update book stock
	book.Stock++
	if err := s.bookRepo.UpdateBook(ctx, tx, book); err != nil {
		return err
	}

	return s.txRepo.Commit(tx)
}

func (s *borrowingRecordService) ListBorrowingRecordsByBookID(ctx context.Context, bookID uuid.UUID) ([]models.BorrowingRecord, error) {
	return s.repo.ListBorrowingRecordsByBookID(ctx, bookID)
}
