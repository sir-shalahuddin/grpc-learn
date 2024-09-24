package auth

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid or expired token")

func ValidateToken(tokenStr, jwtSecret string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and return the secret
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("[ValidateToken] Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Printf("[ValidateToken] Parsing Token : %v", err)
		return uuid.UUID{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			log.Printf("[ValidateToken] Invalid user ID in token claims")
			return uuid.UUID{}, errors.New("invalid user ID in token claims")
		}

		id, err := uuid.Parse(userID)
		if err != nil {
			log.Printf("[ValidateToken] Error parsing uuid : %v", err)
			return uuid.UUID{}, err
		}

		return id, nil
	}

	return uuid.UUID{}, ErrInvalidToken
}
