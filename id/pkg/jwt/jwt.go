package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func NewJWT(privateKey, publicKey []byte) *JWT {
	return &JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j *JWT) GenerateToken(userID int64, login, name string, expiresIn time.Duration) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("generate: parse key: %w", err)
	}

	claims := Claims{
		Login: login,
		Name:  name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (j *JWT) ValidateToken(token string) (bool, Claims) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return false, Claims{}
	}

	claims := Claims{}

	t, err := jwt.ParseWithClaims(token, &claims, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
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
