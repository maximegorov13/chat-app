package service_test

import (
	"context"
	"log"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"

	"github.com/maximegorov13/chat-app/id/internal/rediskeys"
	"github.com/maximegorov13/chat-app/id/pkg/jwt"

	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/auth"
	authredis "github.com/maximegorov13/chat-app/id/internal/auth/repository/redis"
	authservice "github.com/maximegorov13/chat-app/id/internal/auth/service"
	"github.com/maximegorov13/chat-app/id/internal/keyreader"
	"github.com/maximegorov13/chat-app/id/internal/storage/pg"
	"github.com/maximegorov13/chat-app/id/internal/storage/redis"
	"github.com/maximegorov13/chat-app/id/internal/user"
	userpg "github.com/maximegorov13/chat-app/id/internal/user/repository/pg"
	userservice "github.com/maximegorov13/chat-app/id/internal/user/service"
)

func setupTest(t *testing.T) (auth.AuthService, user.UserService, auth.TokenRepository, func(userID int64), func(token string)) {
	conf, err := configs.Load("../../../.env")
	require.NoError(t, err)

	pgClient, err := pg.NewPostgres(conf)
	require.NoError(t, err)

	ctx := context.Background()

	redisClient, err := redis.NewRedis(ctx, conf)
	require.NoError(t, err)

	keyReader := keyreader.NewKeyReader("../../../secrets")
	privateKey, err := keyReader.ReadPrivateKey(conf.Auth.PostfixKeyAuth)
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err := keyReader.ReadPublicKey(conf.Auth.PostfixKeyAuth)
	if err != nil {
		log.Fatal(err)
	}

	jwtMaker := jwt.NewJWT(privateKey, publicKey)

	userRepo := userpg.NewUserRepository(pgClient)
	tokenRepo := authredis.NewTokenRepository(redisClient)

	cleanupUser := func(userID int64) {
		query, args, err := pgClient.Sb.
			Delete("users").
			Where(squirrel.Eq{
				"id": userID,
			}).
			ToSql()
		if err != nil {
			t.Logf("cleanup query build error: %v", err)
		}

		_, err = pgClient.Sqlx.ExecContext(context.Background(), query, args...)
		if err != nil {
			t.Logf("cleanup exec error: %v", err)
		}
	}

	cleanupInvalidToken := func(token string) {
		err := redisClient.Del(ctx, rediskeys.InvalidTokenKey(token))
		if err != nil {
			t.Logf("cleanup exec error: %v", err)
		}
	}

	authService := authservice.NewAuthService(authservice.AuthServiceDeps{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
		JWT:       jwtMaker,
	})
	userService := userservice.NewUserService(userservice.UserServiceDeps{
		UserRepo: userRepo,
	})

	return authService, userService, tokenRepo, cleanupUser, cleanupInvalidToken
}

func TestAuthService_Login(t *testing.T) {
	authService, userService, _, cleanupUser, _ := setupTest(t)

	ctx := context.Background()

	t.Run("successful login", func(t *testing.T) {
		registerReq := &user.RegisterRequest{
			Login:    "testuser",
			Name:     "Test User",
			Password: "12345678",
		}

		registeredUser, err := userService.Register(ctx, registerReq)
		t.Cleanup(func() {
			cleanupUser(registeredUser.ID)
		})
		require.NoError(t, err)

		loginReq := &auth.LoginRequest{
			Login:    registerReq.Login,
			Password: registerReq.Password,
		}

		token, err := authService.Login(ctx, loginReq)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("invalid credentials wrong password", func(t *testing.T) {
		registerReq := &user.RegisterRequest{
			Login:    "testuser",
			Name:     "Test User",
			Password: "12345678",
		}

		registeredUser, err := userService.Register(ctx, registerReq)
		t.Cleanup(func() {
			cleanupUser(registeredUser.ID)
		})
		require.NoError(t, err)

		loginReq := &auth.LoginRequest{
			Login:    registerReq.Login,
			Password: "wrongpassword",
		}

		_, err = authService.Login(ctx, loginReq)
		require.Error(t, err)
		require.ErrorIs(t, err, apperrors.ErrInvalidCredentials)
	})

	t.Run("invalid credentials not found user", func(t *testing.T) {
		loginReq := &auth.LoginRequest{
			Login:    "testuser",
			Password: "12345678",
		}

		_, err := authService.Login(ctx, loginReq)
		require.Error(t, err)
		require.ErrorIs(t, err, apperrors.ErrInvalidCredentials)
	})
}

func TestAuthService_Logout(t *testing.T) {
	authService, userService, tokenRepo, cleanupUser, cleanupInvalidToken := setupTest(t)

	ctx := context.Background()

	t.Run("successful logout", func(t *testing.T) {
		registerReq := &user.RegisterRequest{
			Login:    "testuser",
			Name:     "Test User",
			Password: "12345678",
		}

		registeredUser, err := userService.Register(ctx, registerReq)
		t.Cleanup(func() {
			cleanupUser(registeredUser.ID)
		})
		require.NoError(t, err)

		loginReq := &auth.LoginRequest{
			Login:    registerReq.Login,
			Password: registerReq.Password,
		}

		token, err := authService.Login(ctx, loginReq)
		t.Cleanup(func() {
			cleanupInvalidToken(token)
		})
		require.NoError(t, err)
		require.NotEmpty(t, token)

		err = authService.Logout(ctx, token)
		require.NoError(t, err)

		invalid, err := tokenRepo.IsTokenInvalid(ctx, token)
		require.NoError(t, err)
		require.True(t, invalid)
	})

	t.Run("logout with empty token", func(t *testing.T) {
		err := authService.Logout(ctx, "")
		require.Error(t, err)
		require.ErrorIs(t, err, apperrors.ErrUnauthorized)
	})
}
