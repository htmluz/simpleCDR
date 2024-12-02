package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("chavemtsecreta")

func GenerateAccessToken(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix()
	return token.SignedString(jwtSecret)
}
