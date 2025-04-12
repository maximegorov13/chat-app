package tokencheck_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maximegorov13/chat-app/id/pkg/tokencheck"

	"github.com/maximegorov13/chat-app/id/internal/apperrors"
	"github.com/maximegorov13/chat-app/id/internal/auth"
	"github.com/maximegorov13/chat-app/id/internal/res"
)

func TestClient_IsTokenInvalid(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		expectedToken := "token"
		mockResponse := res.Response[auth.IsTokenInvalidResponse]{
			Data: auth.IsTokenInvalidResponse{
				Invalid: false,
			},
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/auth/is-token-invalid", r.URL.Path)
			require.Equal(t, expectedToken, r.URL.Query().Get("token"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			require.NoError(t, err)
		}))
		defer ts.Close()

		client := tokencheck.NewClient(tokencheck.Config{
			ServiceURL: ts.URL,
		})

		resp, err := client.IsTokenInvalid(context.Background(), expectedToken)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.False(t, resp.Data.Invalid)
	})

	t.Run("invalid token", func(t *testing.T) {
		expectedToken := "token"
		mockResponse := res.Response[auth.IsTokenInvalidResponse]{
			Data: auth.IsTokenInvalidResponse{
				Invalid: true,
			},
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/auth/is-token-invalid", r.URL.Path)
			require.Equal(t, expectedToken, r.URL.Query().Get("token"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			require.NoError(t, err)
		}))
		defer ts.Close()

		client := tokencheck.NewClient(tokencheck.Config{
			ServiceURL: ts.URL,
		})

		resp, err := client.IsTokenInvalid(context.Background(), expectedToken)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.True(t, resp.Data.Invalid)
	})

	t.Run("empty token", func(t *testing.T) {
		expectedToken := ""
		mockResponse := res.Response[auth.IsTokenInvalidResponse]{
			Error: &res.ErrorResponse{
				Code:    apperrors.ErrBadRequest.Code,
				Message: apperrors.ErrBadRequest.Message,
			},
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/auth/is-token-invalid", r.URL.Path)
			require.Equal(t, expectedToken, r.URL.Query().Get(""))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(apperrors.ErrBadRequest.Code)
			err := json.NewEncoder(w).Encode(mockResponse)
			require.NoError(t, err)
		}))
		defer ts.Close()

		client := tokencheck.NewClient(tokencheck.Config{
			ServiceURL: ts.URL,
		})

		resp, err := client.IsTokenInvalid(context.Background(), expectedToken)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, apperrors.ErrBadRequest.Code, resp.Error.Code)
		require.Equal(t, apperrors.ErrBadRequest.Message, resp.Error.Message)
	})

	t.Run("invalid URL", func(t *testing.T) {
		client := tokencheck.NewClient(tokencheck.Config{
			ServiceURL: "http://invalid-url:1234",
		})

		_, err := client.IsTokenInvalid(context.Background(), "token")
		require.Error(t, err)
		var urlErr *url.Error
		require.True(t, errors.As(err, &urlErr))
	})
}
