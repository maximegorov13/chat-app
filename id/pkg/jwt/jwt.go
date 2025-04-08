/*
Package jwt provides JWT (JSON Web Token) generation and validation functionality.

It supports RSA-based signing and verification of tokens with custom claims.
The package is built on top of github.com/golang-jwt/jwt/v5 library.
*/
package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents custom JWT claims along with standard registered claims.
// It includes user login and name in addition to standard JWT fields.
type Claims struct {
	Login string `json:"login"` // User login identifier
	Name  string `json:"name"`  // User display name
	jwt.RegisteredClaims
}

// JWT provides methods for token generation, validation and inspection.
// It requires RSA private and public keys for cryptographic operations.
type JWT struct {
	privateKey []byte // RSA private key in PEM format
	publicKey  []byte // RSA public key in PEM format
}

// NewJWT creates a new JWT instance with provided RSA keys.
// Both privateKey and publicKey should be in PEM format.
func NewJWT(privateKey, publicKey []byte) *JWT {
	return &JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// GenerateToken creates a new JWT token for the specified user.
//
// Parameters:
//   - userID: unique identifier of the user (will be set as 'sub' claim)
//   - login: user login identifier
//   - name: user display name
//   - expiresIn: duration until token expiration
//
// Returns:
//   - signed JWT token string
//   - error if key parsing or signing fails
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

// ValidateToken checks if the token is properly signed and returns its claims.
//
// Parameters:
//   - token: JWT token string to validate
//
// Returns:
//   - bool indicating if token is valid
//   - Claims struct populated with token claims (only valid if first return value is true)
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

// ExtractUserID retrieves user ID from the token's subject claim.
//
// Parameters:
//   - token: JWT token string to inspect
//
// Returns:
//   - user ID as string (from 'sub' claim)
//   - error if token is invalid or doesn't contain subject claim
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

// IsTokenExpired checks if the token has expired.
//
// Parameters:
//   - token: JWT token string to check
//
// Returns:
//   - true if token is expired or invalid, false if still valid
func (j *JWT) IsTokenExpired(token string) bool {
	_, claims := j.ValidateToken(token)
	if claims.ExpiresAt == nil {
		return true
	}

	return claims.ExpiresAt.Before(time.Now())
}
