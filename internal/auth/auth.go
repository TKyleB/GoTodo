package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	TokenSecret                string
	TokenExpirationTime        time.Duration
	RefreshTokenExpirationTime time.Duration
	Issuer                     string
}

type User struct {
	ID        uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (a *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	hashedPassword := string(hash)
	return hashedPassword, nil

}

func (a *AuthService) CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a *AuthService) MakeJWT(userID uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    a.Issuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(a.TokenExpirationTime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.TokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (a *AuthService) ValidateJWT(tokenString string) (uuid.UUID, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.TokenSecret), nil
	})
	if err != nil || !parsedToken.Valid {
		return uuid.UUID{}, err
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.UUID{}, errors.New("invalid Claims Type")
	}

	ID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, err
	}
	return ID, nil

}

func (a *AuthService) GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header on request")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil

}

func (a *AuthService) MakeRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	return hex.EncodeToString(tokenBytes), nil
}
