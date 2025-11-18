package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/ryuyb/fusion/internal/core/command"
	serviceMocks "github.com/ryuyb/fusion/internal/core/port/service"
	"github.com/ryuyb/fusion/internal/infrastructure/http/dto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthController_Register(t *testing.T) {
	app := fiber.New()
	authService := serviceMocks.NewMockAuthService(t)
	controller := NewAuthController(authService, nil)
	app.Post("/auth/register", controller.Register)

	reqBody := `{"username":"neo","email":"neo@example.com","password":"supersecret","confirm_password":"supersecret"}`
	authService.EXPECT().
		Register(mock.Anything, command.RegisterCommand{
			Username: "neo",
			Email:    "neo@example.com",
			Password: "supersecret",
		}).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestAuthController_Login(t *testing.T) {
	app := fiber.New()
	authService := serviceMocks.NewMockAuthService(t)
	controller := NewAuthController(authService, nil)
	app.Post("/auth/login", controller.Login)

	reqBody := `{"username":"neo","password":"supersecret"}`
	expiresAt := time.Now().Add(time.Hour).UTC()
	authService.EXPECT().
		Login(mock.Anything, command.LoginCommand{
			Username: "neo",
			Password: "supersecret",
		}).
		Return("token-123", expiresAt, nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	var loginResp dto.LoginResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&loginResp))
	require.Equal(t, "token-123", loginResp.Token)
	require.True(t, loginResp.ExpireAt.Equal(expiresAt))
}
