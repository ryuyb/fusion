package bark

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/external"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSendValidationErrors(t *testing.T) {
	t.Parallel()
	provider := NewProvider(zap.NewNop())
	ctx := context.Background()
	channel := &domain.NotificationChannel{Config: map[string]any{}}

	err := provider.Send(ctx, channel, &mockNotificationData)
	require.Error(t, err)
	appErr := errors2.GetAppError(err)
	require.Equal(t, errors2.ErrCodeBadRequest, appErr.Code)

	err = provider.Send(ctx, &domain.NotificationChannel{Config: map[string]any{
		"device_key": "abc",
		"level":      "invalid",
	}}, &mockNotificationData)
	require.Equal(t, errors2.ErrCodeBadRequest, errors2.GetAppError(err).Code)

	err = provider.Send(ctx, &domain.NotificationChannel{Config: map[string]any{
		"device_key": "abc",
		"volume":     11,
	}}, &mockNotificationData)
	require.Equal(t, errors2.ErrCodeBadRequest, errors2.GetAppError(err).Code)

	err = provider.Send(ctx, &domain.NotificationChannel{Config: map[string]any{
		"device_key": "abc",
		"icon":       "://bad",
	}}, &mockNotificationData)
	require.Equal(t, errors2.ErrCodeBadRequest, errors2.GetAppError(err).Code)
}

func TestSendHTTPFailure(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("fail"))
	}))
	t.Cleanup(server.Close)

	provider := NewProvider(zap.NewNop())

	err := provider.Send(context.Background(), &domain.NotificationChannel{Config: map[string]any{
		"device_key": "abc",
		"url":        server.URL,
	}}, &mockNotificationData)

	require.Error(t, err)
	appErr := errors2.GetAppError(err)
	require.Equal(t, errors2.ErrCodeBadRequest, appErr.Code)
	require.Equal(t, 500, appErr.Details["status"])
}

func TestSendBarkErrorCode(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    400,
			"message": "token invalid",
		})
	}))
	t.Cleanup(server.Close)

	provider := NewProvider(zap.NewNop())

	err := provider.Send(context.Background(), &domain.NotificationChannel{Config: map[string]any{
		"device_key": "abc",
		"url":        server.URL,
	}}, &mockNotificationData)

	require.Error(t, err)
	appErr := errors2.GetAppError(err)
	require.Equal(t, errors2.ErrCodeBadRequest, appErr.Code)
	require.Equal(t, 400, appErr.Details["code"])
	require.Equal(t, "token invalid", appErr.Details["message"])
}

func TestSendSuccessWithOptionalFields(t *testing.T) {
	t.Parallel()
	received := BarkRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		_ = json.NewDecoder(r.Body).Decode(&received)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    200,
			"message": "ok",
		})
	}))
	t.Cleanup(server.Close)

	provider := NewProvider(zap.NewNop())

	cfg := map[string]any{
		"device_key": "abc",
		"url":        server.URL,
		"subtitle":   "sub",
		"level":      "active",
		"volume":     5,
		"badge":      1,
		"sound":      "ding",
		"group":      "g",
		"icon":       "https://example.com/icon.png",
		"action":     "open",
		"link":       "https://example.com",
	}

	err := provider.Send(context.Background(), &domain.NotificationChannel{Config: cfg}, &mockNotificationData)
	require.NoError(t, err)

	require.Equal(t, "sub", received.Subtitle)
	require.Equal(t, "active", received.Level)
	require.Equal(t, 5, received.Volume)
	require.Equal(t, 1, received.Badge)
	require.Equal(t, "ding", received.Sound)
	require.Equal(t, "g", received.Group)
	require.Equal(t, "https://example.com/icon.png", received.Icon)
	require.Equal(t, "open", received.Action)
	require.Equal(t, "https://example.com", received.Url)
}

// Integration-style test hitting real Bark API; requires network and env opt-in.
func TestSendBarkLiveInvalidDeviceKey(t *testing.T) {
	if testing.Short() || os.Getenv("BARK_LIVE_TEST") != "1" {
		t.Skip("live Bark test skipped; set BARK_LIVE_TEST=1 to run")
	}

	provider := NewProvider(zap.NewNop())
	err := provider.Send(
		context.Background(),
		&domain.NotificationChannel{Config: map[string]any{
			"device_key": "invalid-device-key",
			// use default Bark endpoint
		}},
		&mockNotificationData,
	)

	require.Error(t, err)
	require.Equal(t, errors2.ErrCodeBadRequest, errors2.GetAppError(err).Code)
}

var mockNotificationData = external.NotificationData{
	Title:   "Test",
	Content: "test content",
}
