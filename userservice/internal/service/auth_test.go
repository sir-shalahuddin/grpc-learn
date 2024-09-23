package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/models"
	"github.com/sir-shalahuddin/grpc-learn/userservice/pkg/auth"
)

// MockAuthRepository adalah implementasi mock dari AuthRepository.
type MockAuthRepository struct {
	CreateUserFunc     func(ctx context.Context, user *models.User) error
	GetUserByIDFunc    func(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmailFunc func(ctx context.Context, email string) (*models.User, error)
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	return m.CreateUserFunc(ctx, user)
}

func (m *MockAuthRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return m.GetUserByIDFunc(ctx, userID)
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.GetUserByEmailFunc(ctx, email)
}

// Test Register: Email sudah terdaftar
func TestRegister_EmailAlreadyTaken(t *testing.T) {
	mockRepo := &MockAuthRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{Email: "test@example.com"}, nil
		},
	}
	authService := NewAuthService(mockRepo, "jwt-secret")

	err := authService.Register(context.Background(), dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})

	if err != ErrDuplicateEmail {
		t.Errorf("expected ErrDuplicateEmail, got %v", err)
	}
}

// Test Register: Berhasil mendaftar
func TestRegister_Success(t *testing.T) {
	mockRepo := &MockAuthRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return nil, nil
		},
		CreateUserFunc: func(ctx context.Context, user *models.User) error {
			return nil
		},
	}
	authService := NewAuthService(mockRepo, "jwt-secret")

	err := authService.Register(context.Background(), dto.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		Name:     "New User",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Test Login: Kredensial tidak valid
func TestLogin_InvalidCredentials(t *testing.T) {
	mockRepo := &MockAuthRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return nil, nil
		},
	}
	authService := NewAuthService(mockRepo, "jwt-secret")

	_, err := authService.Login(context.Background(), dto.LoginRequest{
		Email:    "wrong@example.com",
		Password: "password123",
	})

	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

// Test Login: Berhasil login
func TestLogin_Success(t *testing.T) {
	// Hash password for comparison
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo := &MockAuthRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{
				Email:    "user@example.com",
				Password: string(hashedPassword),
				UserID:   uuid.New(),
			}, nil
		},
	}
	authService := NewAuthService(mockRepo, "jwt-secret")

	resp, err := authService.Login(context.Background(), dto.LoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Error("expected valid access and refresh tokens")
	}
}

// Test Refresh Token: Token tidak valid
func TestRefreshToken_InvalidToken(t *testing.T) {
	mockRepo := &MockAuthRepository{}
	authService := NewAuthService(mockRepo, "jwt-secret")

	_, err := authService.RefreshToken(context.Background(), dto.RefreshTokenRequest{
		RefreshToken: "invalid-token",
	})

	if err != auth.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

// Test Refresh Token: Berhasil refresh token
func TestRefreshToken_Success(t *testing.T) {
	userID := uuid.New()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(RefreshTokenExpiry).Unix(),
	})
	refreshToken, _ := token.SignedString([]byte("jwt-secret"))

	mockRepo := &MockAuthRepository{}
	authService := NewAuthService(mockRepo, "jwt-secret")

	resp, err := authService.RefreshToken(context.Background(), dto.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resp.AccessToken == "" {
		t.Error("expected a new access token")
	}
}
