package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/dto"
	"github.com/sir-shalahuddin/grpc-learn/userservice/internal/models"
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
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	registeredUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if registeredUser != nil {
		return ErrDuplicateEmail
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[Service - Register] Error hashing password: %v", err)
		return err
	}

	user := models.User{
		Email:    req.Email,
		Password: string(passwordHash),
		Name:     req.Name,
	}

	if err := s.repo.CreateUser(ctx, &user); err != nil {
		return err
	}

	return nil
}

// Login user with email and password, returning access and refresh tokens.
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	// Retrieve the user
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if user == nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	// Check the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Printf("[Service - Login] Error comparing password: %v", err)
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	// Generate Access Token
	accessToken, err := s.generateToken(user.UserID, "access")
	if err != nil {
		log.Printf("[Service - Login] Error generate token: %v", err)
		return dto.LoginResponse{}, err
	}
	// Generate Refresh Token
	refreshToken, err := s.generateToken(user.UserID, "refresh")
	if err != nil {
		log.Printf("[Service - Login] Error generate token: %v", err)
		return dto.LoginResponse{}, err
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

// RefreshToken generates a new access token using the provided refresh token.
func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error) {
	// Parse and validate the refresh token
	userID, err := s.ValidateToken(req.RefreshToken)
	if err != nil {
		return dto.RefreshTokenResponse{}, ErrInvalidToken
	}

	// Generate a new access token
	newAccessToken, err := s.generateToken(userID, "access")
	if err != nil {
		log.Printf("[Service - RefreshToken] Error generate token: %v", err)
		return dto.RefreshTokenResponse{}, err
	}

	return dto.RefreshTokenResponse{AccessToken: newAccessToken}, nil
}

// GetUserByID retrieves a user by their ID.
func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
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
		return "", err
	}

	return signedToken, nil
}

// ValidateToken parses and validates the JWT token, returning the userID if valid.
func (s *authService) ValidateToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and return the secret
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("[Service - ValidateToken] Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		log.Printf("[Service - ValidateToken] Parsing Token : %v", err)
		return uuid.UUID{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			log.Printf("[Service - ValidateToken] Invalid user ID in token claims")
			return uuid.UUID{}, errors.New("invalid user ID in token claims")
		}

		id, err := uuid.Parse(userID)
		if err != nil {
			log.Printf("[Service - ValidateToken] Error parsing uuid : %v", err)
			return uuid.UUID{}, err
		}

		return id, nil
	}

	return uuid.UUID{}, ErrInvalidToken
}
