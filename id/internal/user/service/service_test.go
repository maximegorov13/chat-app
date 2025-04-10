package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/maximegorov13/chat-app/id/configs"
	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/storage/pg"
	"github.com/maximegorov13/chat-app/id/internal/user"
	userpg "github.com/maximegorov13/chat-app/id/internal/user/repository/pg"
	userservice "github.com/maximegorov13/chat-app/id/internal/user/service"
)

type testDependencies struct {
	userService user.UserService
	userRepo    user.UserRepository
	cleanupUser func(userID int64)
}

func getUniqueLogin() string {
	return fmt.Sprintf("user-%s", uuid.New())
}

func setupTest(t testing.TB) *testDependencies {
	t.Helper()

	conf, err := configs.Load("../../../.env")
	require.NoError(t, err)

	pgClient, err := pg.NewPostgres(conf)
	require.NoError(t, err)

	userRepo := userpg.NewUserRepository(pgClient)

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

	return &testDependencies{
		userService: userservice.NewUserService(userservice.UserServiceDeps{
			UserRepo: userRepo,
		}),
		userRepo:    userRepo,
		cleanupUser: cleanupUser,
	}
}

func TestUserService_Register(t *testing.T) {
	deps := setupTest(t)

	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		uniqueLogin := getUniqueLogin()
		registerReq := &user.RegisterRequest{
			Login:    uniqueLogin,
			Name:     "Test User",
			Password: "12345678",
		}

		u, err := deps.userService.Register(ctx, registerReq)
		t.Cleanup(func() {
			deps.cleanupUser(u.ID)
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

		dbUser, err := deps.userRepo.FindByLogin(ctx, registerReq.Login)
		require.NoError(t, err)
		require.NotNil(t, dbUser)
		require.Equal(t, u, dbUser)
	})

	t.Run("user already exists", func(t *testing.T) {
		uniqueLogin := getUniqueLogin()
		registerReq := &user.RegisterRequest{
			Login:    uniqueLogin,
			Name:     "Test User",
			Password: "12345678",
		}

		u, err := deps.userService.Register(ctx, registerReq)
		t.Cleanup(func() {
			deps.cleanupUser(u.ID)
		})
		require.NoError(t, err)

		_, err = deps.userService.Register(ctx, registerReq)
		require.Error(t, err)
		require.ErrorIs(t, err, apperrors.ErrUserExists)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	deps := setupTest(t)

	ctx := context.Background()

	t.Run("successful update user", func(t *testing.T) {
		uniqueLogin := getUniqueLogin()
		registerReq := &user.RegisterRequest{
			Login:    uniqueLogin,
			Name:     "Original Name",
			Password: "originalpass",
		}

		registeredUser, err := deps.userService.Register(ctx, registerReq)
		t.Cleanup(func() {
			deps.cleanupUser(registeredUser.ID)
		})
		require.NoError(t, err)

		updateReq := &user.UpdateUserRequest{
			Login:    "updated_login",
			Name:     "Updated Name",
			Password: "newpassword123",
		}
		updatedUser, err := deps.userService.UpdateUser(ctx, registeredUser.ID, updateReq)
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
		uniqueLogin := getUniqueLogin()
		updateReq := &user.UpdateUserRequest{
			Login:    uniqueLogin,
			Name:     "Updated Name",
			Password: "newpassword123",
		}
		_, err := deps.userService.UpdateUser(ctx, 9999999999, updateReq)
		require.Error(t, err)
		require.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}
