package auth

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("veryinsecuresecret")

func GenerateJWT(username string, perms bool) (string, error) {
	claims := jwt.MapClaims{
		"sub":   username,
		"admin": perms,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		slog.Error("Failed to create jwt token", "error", err.Error())
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(token string) (jwt.MapClaims, error) {
	tokenParsed, err := jwt.Parse(token, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Error("Unexpected signing method: ", "token_method", tkn.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method %v", tkn.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		slog.Error("Failed to parse the token in verify jwt", "error", err.Error())
		return nil, err
	}

	if claims, ok := tokenParsed.Claims.(jwt.MapClaims); ok && tokenParsed.Valid {
		return claims, nil
	}

	slog.Error("Invalid token", "errtkn", token)
	return nil, fmt.Errorf("invalid token")
}
