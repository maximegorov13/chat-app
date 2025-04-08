package service_test

import (
	"context"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/storage/pg"
	"github.com/maximegorov13/chat-app/id/internal/user"
	userpg "github.com/maximegorov13/chat-app/id/internal/user/repository/pg"
	userservice "github.com/maximegorov13/chat-app/id/internal/user/service"
)

func setupTest(t testing.TB) (user.UserService, user.UserRepository, func(userID int64)) {
	t.Helper()

	conf, err := configs.Load("../../../.env")
	require.NoError(t, err)

	pgClient, err := pg.NewPostgres(conf)
	require.NoError(t, err)

	userRepo := userpg.NewUserRepository(pgClient)
	userService := userservice.NewUserService(userservice.UserServiceDeps{
		UserRepo: userRepo,
	})

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

	return userService, userRepo, cleanupUser
}

func TestUserService_Register(t *testing.T) {
	userService, repo, cleanupUser := setupTest(t)

	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		registerReq := &user.RegisterRequest{
			Login:    "testuser",
			Name:     "Test User",
			Password: "12345678",
		}

		u, err := userService.Register(ctx, registerReq)
		t.Cleanup(func() {
			cleanupUser(u.ID)
		})
		require.NoError(t, err)
		require.NotNil(t, u)
		require.NotZero(t, u.ID)
		require.Equal(t, registerReq.Login, u.Login)
		require.Equal(t, registerReq.Name, u.Name)
		require.NotEmpty(t, u.Password)

		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(registerReq.Password))
		require.NoError(t, err)
		require.False(t, u.CreatedAt.IsZero())
		require.False(t, u.UpdatedAt.IsZero())

		dbUser, err := repo.FindByLogin(ctx, registerReq.Login)
		require.NoError(t, err)
		require.NotNil(t, dbUser)
		require.Equal(t, u, dbUser)
	})

	t.Run("user already exists", func(t *testing.T) {
		registerReq := &user.RegisterRequest{
			Login:    "testuser",
			Name:     "Test User",
			Password: "12345678",
		}

		u, err := userService.Register(ctx, registerReq)
		t.Cleanup(func() {
			cleanupUser(u.ID)
		})
		require.NoError(t, err)

		_, err = userService.Register(ctx, registerReq)
		require.Error(t, err)
		require.ErrorIs(t, err, apperrors.ErrUserExists)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	userService, _, cleanupUser := setupTest(t)

	ctx := context.Background()

	t.Run("successful update user", func(t *testing.T) {
		registerReq := &user.RegisterRequest{
			Login:    "user_to_update",
			Name:     "Original Name",
			Password: "originalpass",
		}

		registeredUser, err := userService.Register(ctx, registerReq)
		t.Cleanup(func() {
			cleanupUser(registeredUser.ID)
		})
		require.NoError(t, err)

		updateReq := &user.UpdateUserRequest{
			Login:    "updated_login",
			Name:     "Updated Name",
			Password: "newpassword123",
		}
		updatedUser, err := userService.UpdateUser(ctx, registeredUser.ID, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updatedUser)
		require.Equal(t, updateReq.Login, updatedUser.Login)
		require.Equal(t, updateReq.Name, updatedUser.Name)

		err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(updateReq.Password))
		require.NoError(t, err)

		err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(registerReq.Password))
		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		updateReq := &user.UpdateUserRequest{
			Login:    "updated_login",
			Name:     "Updated Name",
			Password: "newpassword123",
		}
		_, err := userService.UpdateUser(ctx, 9999999999, updateReq)
		require.Error(t, err)
		require.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}
