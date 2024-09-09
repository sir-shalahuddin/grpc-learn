package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenExpiry  = time.Minute * 60   // 15 minutes for access tokens
	RefreshTokenExpiry = time.Hour * 24 * 7 // 7 days for refresh tokens
)

var (
	ErrDuplicateEmail     = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrInsufficientAccess = errors.New("access forbidden: insufficient permissions")
)

type authService struct {
	userRepo  UserRepository
	jwtSecret string
}

func NewAuthService(userRepo UserRepository, jwtSecret string) *authService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, email, password, name string) error {

	registeredUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if registeredUser != nil {
		return ErrDuplicateEmail
	}
	if err != nil {
		return err
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    email,
		Password: string(passwordHash),
		Name:     name,
	}

	err = s.userRepo.CreateUser(context.Background(), &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(ctx context.Context, email, password string) (map[string]string, error) {
	// Retrieve the user
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// Generate Refresh Token
	refreshToken, err := s.generateToken(user.UserID, "refresh")
	if err != nil {
		return nil, err
	}

	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return tokens, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Parse and validate the refresh token
	log.Printf(refreshToken)
	userID, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Generate a new access token
	newAccessToken, err := s.generateToken(userID, "access")
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) generateToken(userID uuid.UUID, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // Default to access token expiry
	}

	if tokenType == "refresh" {
		// Set longer expiration for refresh token
		claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // Example: 7 days
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *authService) ValidateToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and return the secret
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			return uuid.UUID{}, errors.New("invalid user ID in token claims")
		}

		id, err := uuid.Parse(userID)
		if err != nil {
			return uuid.UUID{}, errors.New("failed to parse user id")
		}

		return id, nil
	}

	return uuid.UUID{}, ErrInvalidToken
}

func (s *authService) HasAccess(userRole string, allowedRoles []string) bool {
	if len(allowedRoles) == 0 {
		allowedRoles = append(allowedRoles, "user")
	}
	if userRole == "super admin" {
		return true
	}

	for _, role := range allowedRoles {
		if role == userRole {
			return true
		}
	}
	return false
}
