package utils

import (
	"blog/internal/config"
	"blog/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user *models.User, cfg *config.Config) (string, error) {
	claims := &Claims{
		ID:    user.ID.String(),
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(cfg.HTTPServer.JWTSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
