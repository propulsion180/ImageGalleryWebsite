package auth

import (
	"fmt"
	"log"
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
		log.Println("failed to create jwt token: ", err.Error())
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(token string) (jwt.MapClaims, error) {
	tokenParsed, err := jwt.Parse(token, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("unexpected signing method: ", tkn.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method %v", tkn.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		log.Println("failed to parse the token in verify jwt")
		return nil, err
	}

	if claims, ok := tokenParsed.Claims.(jwt.MapClaims); ok && tokenParsed.Valid {
		return claims, nil
	}

	log.Println("invalid token: ", token)
	return nil, fmt.Errorf("invalid token")
}
