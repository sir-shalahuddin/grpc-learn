package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/bookservice/models"
)

var (
	ErrBorrowingRecordNotFound = errors.New("borrowing record not found")
	ErrBookUnavailable         = errors.New("failed to process due to 0 stock")
)

type BorrowingRecordRepository interface {
	CreateBorrowingRecord(ctx context.Context, tx *sql.Tx, record *models.BorrowingRecord) error
	GetBorrowingRecordByID(ctx context.Context, id uuid.UUID) (*models.BorrowingRecord, error)
	UpdateBorrowingRecord(ctx context.Context, tx *sql.Tx, record *models.BorrowingRecord) error
	DeleteBorrowingRecord(ctx context.Context, id uuid.UUID) error
	ListBorrowingRecords(ctx context.Context, userID uuid.UUID, queries map[string]string) ([]models.BorrowingRecord, error)
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

func (s *borrowingRecordService) BorrowBook(ctx context.Context, req dto.BorrowBookRequest, bookID, userID uuid.UUID) error {
	book, err := s.bookRepo.GetBookByID(ctx, bookID)
	if err != nil {
		return err
	}
	if book == nil {
		return ErrBookNotFound
	}
	if book.Stock <= 0 {
		return ErrBookUnavailable
	}

	record := &models.BorrowingRecord{
		Book: models.Book{
			ID: bookID,
		},
		UserID:  userID,
		DueDate: req.DueDate,
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

func (s *borrowingRecordService) ReturnBook(ctx context.Context, bookID, recordID uuid.UUID) error {
	record, err := s.repo.GetBorrowingRecordByID(ctx, recordID)
	if err != nil {
		return err
	}
	if record == nil {
		return ErrBorrowingRecordNotFound
	}

	book, err := s.bookRepo.GetBookByID(ctx, record.Book.ID)
	if err != nil {
		return err
	}
	if book == nil {
		return ErrBookNotFound
	}

	// Update record
	returnAt := time.Now()
	record.ReturnedAt = &returnAt
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

func (s *borrowingRecordService) ListBorrowingRecords(ctx context.Context, queries map[string]string, userID uuid.UUID) ([]models.BorrowingRecord, error) {

	record, err := s.repo.ListBorrowingRecords(ctx, userID, queries)
	if err != nil {
		return nil, err
	}

	if len(record) == 0 {
		return nil, nil
	}

	return record, nil
}
