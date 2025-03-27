package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type Claims struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

type JWT struct {
	Secret string
}

func New(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) GenerateToken(userID int64, login, name string, expiresIn time.Duration) (string, error) {
	claims := Claims{
		Login: login,
		Name:  name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString([]byte(j.Secret))
}

func (j *JWT) ValidateToken(token string) (bool, Claims) {
	claims := Claims{}
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.Secret), nil
	})
	if err != nil || !t.Valid {
		return false, Claims{}
	}

	return true, claims
}

func (j *JWT) ExtractUserID(token string) (string, error) {
	valid, claims := j.ValidateToken(token)
	if !valid {
		return "", fmt.Errorf("invalid token")
	}

	if claims.Subject == "" {
		return "", fmt.Errorf("token does not contain subject (sub)")
	}

	return claims.Subject, nil
}

func (j *JWT) IsTokenExpired(token string) bool {
	_, claims := j.ValidateToken(token)
	if claims.ExpiresAt == nil {
		return false
	}

	return claims.ExpiresAt.Before(time.Now())
}
