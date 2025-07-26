package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/hexagonal/adapters/drivers/rest/dto"
)

type mockService struct {
	registerFunc func(ctx context.Context, input dto.RegisterDTO) error
	loginFunc    func(ctx context.Context, input dto.LoginDTO) (string, error)
}

func (m *mockService) Register(ctx context.Context, input dto.RegisterDTO) error {
	return m.registerFunc(ctx, input)
}

func (m *mockService) Login(ctx context.Context, input dto.LoginDTO) (string, error) {
	return m.loginFunc(ctx, input)
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockService    *mockService
		body           any
		expectedStatus int
	}{
		{
			name: "success",
			mockService: &mockService{
				registerFunc: func(ctx context.Context, input dto.RegisterDTO) error {
					return nil
				},
			},
			body: dto.RegisterDTO{
				Email:    "john@example.com",
				Password: "secret",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "bad request",
			mockService:    &mockService{},
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal server error",
			mockService: &mockService{
				registerFunc: func(ctx context.Context, input dto.RegisterDTO) error {
					return errors.New("something went wrong")
				},
			},
			body: dto.RegisterDTO{
				Email:    "jane@example.com",
				Password: "123456",
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewAdminHandler(tt.mockService)
			router := gin.New()
			router.POST("/register", handler.Register)

			var reqBody io.Reader
			switch v := tt.body.(type) {
			case string:
				reqBody = bytes.NewBufferString(v)
			default:
				jsonPayload, _ := json.Marshal(v)
				reqBody = bytes.NewBuffer(jsonPayload)
			}

			resp := performRequest(router, "POST", "/register", reqBody)
			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockService    *mockService
		body           any
		expectedStatus int
		expectedToken  string
	}{
		{
			name: "success",
			mockService: &mockService{
				loginFunc: func(ctx context.Context, input dto.LoginDTO) (string, error) {
					return "mock-token", nil
				},
			},
			body: dto.LoginDTO{
				Email:    "john@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedToken:  "mock-token",
		},
		{
			name:           "bad request",
			mockService:    &mockService{},
			body:           "bad json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal server error",
			mockService: &mockService{
				loginFunc: func(ctx context.Context, input dto.LoginDTO) (string, error) {
					return "", errors.New("auth failed")
				},
			},
			body: dto.LoginDTO{
				Email:    "john@example.com",
				Password: "wrong",
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewAdminHandler(tt.mockService)
			router := gin.New()
			router.POST("/login", handler.Login)

			var reqBody io.Reader
			switch v := tt.body.(type) {
			case string:
				reqBody = bytes.NewBufferString(v)
			default:
				jsonPayload, _ := json.Marshal(v)
				reqBody = bytes.NewBuffer(jsonPayload)
			}

			resp := performRequest(router, "POST", "/login", reqBody)
			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var responseBody map[string]string
				json.NewDecoder(resp.Body).Decode(&responseBody)
				if responseBody["token"] != tt.expectedToken {
					t.Errorf("expected token '%s', got '%s'", tt.expectedToken, responseBody["token"])
				}
			}
		})
	}
}
