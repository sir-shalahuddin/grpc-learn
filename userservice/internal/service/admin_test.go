package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/models"
)

// Mock AdminRepository untuk pengujian
type MockAdminRepository struct {
	GetUserByIDFunc    func(ctx context.Context, userID uuid.UUID) (*models.User, error)
	DeleteUserFunc     func(ctx context.Context, userID uuid.UUID) error
	ListUsersFunc      func(ctx context.Context) ([]dto.GetUser, error)
	UpdateUserRolesFunc func(ctx context.Context, userID uuid.UUID, roles string) error
}

func (m *MockAdminRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return m.GetUserByIDFunc(ctx, userID)
}

func (m *MockAdminRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return m.DeleteUserFunc(ctx, userID)
}

func (m *MockAdminRepository) ListUsers(ctx context.Context) ([]dto.GetUser, error) {
	return m.ListUsersFunc(ctx)
}

func (m *MockAdminRepository) UpdateUserRoles(ctx context.Context, userID uuid.UUID, roles string) error {
	return m.UpdateUserRolesFunc(ctx, userID, roles)
}

// Test ListUsers: Berhasil mendapatkan daftar pengguna
func TestListUsers_Success(t *testing.T) {
	mockRepo := &MockAdminRepository{
		ListUsersFunc: func(ctx context.Context) ([]dto.GetUser, error) {
			return []dto.GetUser{
				{UserID: uuid.New(), Name: "User1", Email: "user1@example.com"},
				{UserID: uuid.New(), Name: "User2", Email: "user2@example.com"},
			}, nil
		},
	}
	adminService := NewAdminService(mockRepo)

	users, err := adminService.ListUsers(context.Background())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

// Test UpdateUserRoles: Role berhasil diperbarui
func TestUpdateUserRoles_Success(t *testing.T) {
	mockRepo := &MockAdminRepository{
		GetUserByIDFunc: func(ctx context.Context, userID uuid.UUID) (*models.User, error) {
			return &models.User{
				UserID: userID,
				Role:   "user",
			}, nil
		},
		UpdateUserRolesFunc: func(ctx context.Context, userID uuid.UUID, roles string) error {
			return nil
		},
	}
	adminService := NewAdminService(mockRepo)

	err := adminService.UpdateUserRoles(context.Background(), uuid.New(), dto.UpdateUserRoles{Role: "admin"})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Test UpdateUserRoles: Menghindari promosi menjadi "super admin"
func TestUpdateUserRoles_InsufficientPermissions(t *testing.T) {
	mockRepo := &MockAdminRepository{
		GetUserByIDFunc: func(ctx context.Context, userID uuid.UUID) (*models.User, error) {
			return &models.User{
				UserID: userID,
				Role:   "user",
			}, nil
		},
	}
	adminService := NewAdminService(mockRepo)

	err := adminService.UpdateUserRoles(context.Background(), uuid.New(), dto.UpdateUserRoles{Role: "super admin"})

	if !errors.Is(err, ErrInsufficientPermissions) {
		t.Errorf("expected ErrInsufficientPermissions, got %v", err)
	}
}

// Test DeleteUser: Berhasil menghapus pengguna
func TestDeleteUser_Success(t *testing.T) {
	mockRepo := &MockAdminRepository{
		GetUserByIDFunc: func(ctx context.Context, userID uuid.UUID) (*models.User, error) {
			return &models.User{
				UserID: userID,
				Role:   "user",
			}, nil
		},
		DeleteUserFunc: func(ctx context.Context, userID uuid.UUID) error {
			return nil
		},
	}
	adminService := NewAdminService(mockRepo)

	err := adminService.DeleteUser(context.Background(), uuid.New())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Test DeleteUser: Tidak dapat menghapus pengguna dengan role "super admin"
func TestDeleteUser_InsufficientPermissions(t *testing.T) {
	mockRepo := &MockAdminRepository{
		GetUserByIDFunc: func(ctx context.Context, userID uuid.UUID) (*models.User, error) {
			return &models.User{
				UserID: userID,
				Role:   "super admin",
			}, nil
		},
	}
	adminService := NewAdminService(mockRepo)

	err := adminService.DeleteUser(context.Background(), uuid.New())

	if !errors.Is(err, ErrInsufficientPermissions) {
		t.Errorf("expected ErrInsufficientPermissions, got %v", err)
	}
}

// Test DeleteUser: Pengguna tidak ditemukan
func TestDeleteUser_UserNotFound(t *testing.T) {
	mockRepo := &MockAdminRepository{
		GetUserByIDFunc: func(ctx context.Context, userID uuid.UUID) (*models.User, error) {
			return nil, nil
		},
	}
	adminService := NewAdminService(mockRepo)

	err := adminService.DeleteUser(context.Background(), uuid.New())

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
