package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID  int64 `json:"uid"`
	IsAdmin bool  `json:"is_admin"`
	jwt.RegisteredClaims
}

func NewToken(userID int64, isAdmin bool, ttl time.Duration) (string, error) {
	secret, _ := os.LookupEnv("SecretKey")
	Key := []byte(secret)

	claims := Claims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(Key)
}

func ParseToken(tokenString string) (*Claims, error) {
	sec, _ := os.LookupEnv("SecretKey")
	SecretKey := []byte(sec)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sugning method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
