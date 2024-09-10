package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenExpiry  = time.Minute * 60   // 60 minutes for access tokens
	RefreshTokenExpiry = time.Hour * 24 * 7 // 7 days for refresh tokens
)

var (
	ErrDuplicateEmail     = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrInsufficientAccess = errors.New("access forbidden: insufficient permissions")
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type authService struct {
	repo      AuthRepository
	jwtSecret string
}

func NewAuthService(repo AuthRepository, jwtSecret string) *authService {
	return &authService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Register a new user with email, password, and name.
func (s *authService) Register(ctx context.Context, email, password, name string) error {
	registeredUser, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("service: error checking email existence: %w", err)
	}
	if registeredUser != nil {
		return ErrDuplicateEmail
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("service: error hashing password: %w", err)
	}

	user := models.User{
		Email:    email,
		Password: string(passwordHash),
		Name:     name,
	}

	err = s.repo.CreateUser(ctx, &user)
	if err != nil {
		return fmt.Errorf("service: error creating user: %w", err)
	}

	return nil
}

// Login user with email and password, returning access and refresh tokens.
func (s *authService) Login(ctx context.Context, email, password string) (map[string]string, error) {
	// Retrieve the user
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("service: error retrieving user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate Access Token
	accessToken, err := s.generateToken(user.UserID, "access")
	if err != nil {
		return nil, fmt.Errorf("service: error generating access token: %w", err)
	}

	// Generate Refresh Token
	refreshToken, err := s.generateToken(user.UserID, "refresh")
	if err != nil {
		return nil, fmt.Errorf("service: error generating refresh token: %w", err)
	}

	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return tokens, nil
}

// RefreshToken generates a new access token using the provided refresh token.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Parse and validate the refresh token
	userID, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("service: error validating refresh token: %w", err)
	}

	// Generate a new access token
	newAccessToken, err := s.generateToken(userID, "access")
	if err != nil {
		return "", fmt.Errorf("service: error generating new access token: %w", err)
	}

	return newAccessToken, nil
}

// GetUserByID retrieves a user by their ID.
func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: error retrieving user by ID: %w", err)
	}
	return user, nil
}

// generateToken creates a JWT token with the specified userID and tokenType.
func (s *authService) generateToken(userID uuid.UUID, tokenType string) (string, error) {
	expiry := AccessTokenExpiry
	if tokenType == "refresh" {
		expiry = RefreshTokenExpiry
	}

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(expiry).Unix(),
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("service: error signing token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken parses and validates the JWT token, returning the userID if valid.
func (s *authService) ValidateToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and return the secret
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("service: error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			return uuid.UUID{}, errors.New("invalid user ID in token claims")
		}

		id, err := uuid.Parse(userID)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("service: failed to parse user ID: %w", err)
		}

		return id, nil
	}

	return uuid.UUID{}, ErrInvalidToken
}
