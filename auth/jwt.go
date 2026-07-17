package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, secret string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "chess-puzzle-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ValidateToken(tokenString, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) { //why the callback? because the json is decoded but the libary func is asking what key to use to verify it?
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //the callback gives us control on this because different JWTs may have used different keys.
			return nil, fmt.Errorf("unexepected signing method: %v", token.Header["alg"]) // Don't hand the parser a key—hand it a way to find the right one when the token gives up enough clues.
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}
	return "", errors.New("invalid token structure or expired")
}
