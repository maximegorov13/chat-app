package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/maximegorov13/chat-app/id/pkg/jwt"
)

func generateTestRSAKeys(t testing.TB) (privateKey, publicKey []byte) {
	t.Helper()

	privateKeyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKeyPair.PublicKey)
	require.NoError(t, err)

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKeyPair)

	privateKey = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKey = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKey, publicKey
}

func TestJWT(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	j := jwt.NewJWT(privateKey, publicKey)

	userID := int64(1)
	login := "testuser"
	name := "Test User"
	expiresIn := time.Hour

	t.Run("generate and validate token", func(t *testing.T) {
		token, err := j.GenerateToken(userID, login, name, expiresIn)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		valid, claims := j.ValidateToken(token)
		require.True(t, valid)
		require.Equal(t, login, claims.Login)
		require.Equal(t, name, claims.Name)
		require.Equal(t, strconv.FormatInt(userID, 10), claims.Subject)
		require.True(t, claims.ExpiresAt.After(time.Now()))
	})

	t.Run("invalid token", func(t *testing.T) {
		valid, _ := j.ValidateToken("invalid_token")
		require.False(t, valid)
	})

	t.Run("extract user id", func(t *testing.T) {
		token, err := j.GenerateToken(userID, login, name, expiresIn)
		require.NoError(t, err)

		extractedID, err := j.ExtractUserID(token)
		require.NoError(t, err)
		require.Equal(t, strconv.FormatInt(userID, 10), extractedID)
	})

	t.Run("extract user id from invalid token", func(t *testing.T) {
		_, err := j.ExtractUserID("invalid_token")
		require.Error(t, err)
	})

	t.Run("token expiration", func(t *testing.T) {
		token, err := j.GenerateToken(userID, login, name, -time.Hour)
		require.NoError(t, err)
		require.True(t, j.IsTokenExpired(token))
	})

	t.Run("token not expiration", func(t *testing.T) {
		token, err := j.GenerateToken(userID, login, name, time.Hour)
		require.NoError(t, err)
		require.False(t, j.IsTokenExpired(token))
	})

	t.Run("invalid keys", func(t *testing.T) {
		invalidJWT := jwt.NewJWT([]byte("invalid"), []byte("invalid"))

		_, err := invalidJWT.GenerateToken(userID, login, name, expiresIn)
		require.Error(t, err)

		valid, _ := invalidJWT.ValidateToken("token")
		require.False(t, valid)
	})
}

func ExampleJWT_GenerateToken() {
	privateKey := `-----BEGIN RSA PRIVATE KEY-----...`
	publicKey := `-----BEGIN PUBLIC KEY-----...`

	j := jwt.NewJWT([]byte(privateKey), []byte(publicKey))
	token, err := j.GenerateToken(1, "test", "Test", time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
}

func ExampleJWT_ValidateToken() {
	privateKey := `-----BEGIN RSA PRIVATE KEY-----...`
	publicKey := `-----BEGIN PUBLIC KEY-----...`

	j := jwt.NewJWT([]byte(privateKey), []byte(publicKey))
	token, err := j.GenerateToken(1, "test", "Test", time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	valid, claims := j.ValidateToken(token)
	fmt.Println(valid, claims)
}

func ExampleJWT_ExtractUserID() {
	privateKey := `-----BEGIN RSA PRIVATE KEY-----...`
	publicKey := `-----BEGIN PUBLIC KEY-----...`

	j := jwt.NewJWT([]byte(privateKey), []byte(publicKey))
	token, err := j.GenerateToken(1, "test", "Test", time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	userID, err := j.ExtractUserID(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userID)
}

func ExampleJWT_IsTokenExpired() {
	privateKey := `-----BEGIN RSA PRIVATE KEY-----...`
	publicKey := `-----BEGIN PUBLIC KEY-----...`

	j := jwt.NewJWT([]byte(privateKey), []byte(publicKey))
	token, err := j.GenerateToken(1, "test", "Test", time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	isExpired := j.IsTokenExpired(token)
	fmt.Println(isExpired)
}
