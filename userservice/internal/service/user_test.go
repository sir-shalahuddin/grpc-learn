package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/models"
)

// Mock UserRepository untuk pengujian
type MockUserRepository struct {
	GetUserByIDFunc    func(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUserFunc     func(ctx context.Context, user *models.User) error
	GetUserByEmailFunc func(ctx context.Context, email string) (*models.User, error)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return m.GetUserByIDFunc(ctx, userID)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return m.UpdateUserFunc(ctx, user)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.GetUserByEmailFunc(ctx, email)
}

// Test GetUserByID: User ditemukan
func TestGetUserByID_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetUserByIDFunc: func(ctx context.Context, userID uuid.UUID) (*models.User, error) {
			return &models.User{
				UserID:    userID,
				Name:      "Test User",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}
	userService := NewUserService(mockRepo)

	userID := uuid.New()
	resp, err := userService.GetUserByID(context.Background(), userID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, resp.UserID)
	}
}

// Test GetUserByID: User tidak ditemukan
func TestGetUserByID_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetUserByIDFunc: func(ctx context.Context, userID uuid.UUID) (*models.User, error) {
			return nil, nil
		},
	}
	userService := NewUserService(mockRepo)

	_, err := userService.GetUserByID(context.Background(), uuid.New())

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

// Test UpdateUser: Email sudah digunakan oleh user lain
func TestUpdateUser_EmailAlreadyTaken(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{
				UserID: uuid.New(),
				Email:  email,
			}, nil
		},
	}
	userService := NewUserService(mockRepo)

	err := userService.UpdateUser(context.Background(), uuid.New(), dto.UpdateProfileRequest{
		Email: "test@example.com",
		Name:  "New Name",
	})

	if !errors.Is(err, ErrDuplicateEmail) {
		t.Errorf("expected ErrDuplicateEmail, got %v", err)
	}
}

// Test UpdateUser: Berhasil update
func TestUpdateUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return nil, nil
		},
		UpdateUserFunc: func(ctx context.Context, user *models.User) error {
			return nil
		},
	}
	userService := NewUserService(mockRepo)

	userID := uuid.New()
	err := userService.UpdateUser(context.Background(), userID, dto.UpdateProfileRequest{
		Email: "newemail@example.com",
		Name:  "New Name",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
