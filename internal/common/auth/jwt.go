package auth

import (
	"errors"
	"fmt"
	"github.com/drunkleen/rasta/config"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateJWTToken generates a JWT token based on the provided email and user ID.
//
// Parameter email is the user's email address and userId is the unique identifier of the user.
// Return type is a string representing the generated JWT token and an error object that is returned if the generation fails.
func GenerateJWTToken(email string, userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Second * time.Duration(config.GetJwtExpiry())).Unix(),
	})
	return token.SignedString([]byte(config.GetJwtSecret()))
}

// ValidateJWTToken validates a JWT token and returns the user ID associated with it.
//
// Parameter token is the JWT token to be validated.
// Return type is a string representing the user ID and an error object that is returned if the validation fails.
func ValidateJWTToken(token string) (string, string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
		}
		return []byte(config.GetJwtSecret()), nil
	})
	if err != nil {
		return "", "", fmt.Errorf("token parsing error: %w", err)
	}
	if !parsedToken.Valid {
		return "", "", errors.New("invalid or expired token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}
	return userId, email, nil
}
