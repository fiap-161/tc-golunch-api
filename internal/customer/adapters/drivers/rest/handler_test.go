package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest/dto"
	"github.com/gin-gonic/gin"
)

type mockCustomerService struct {
	createFunc   func(ctx context.Context, dto dto.CreateCustomerDTO) (string, error)
	identifyFunc func(ctx context.Context, CPF string) (string, error)
}

func (m *mockCustomerService) Create(ctx context.Context, dto dto.CreateCustomerDTO) (string, error) {
	return m.createFunc(ctx, dto)
}

func (m *mockCustomerService) Identify(ctx context.Context, CPF string) (string, error) {
	return m.identifyFunc(ctx, CPF)
}

func TestCustomerHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockCreateFunc func(ctx context.Context, dto dto.CreateCustomerDTO) (string, error)
		expectedCode   int
		expectedBody   string
	}{
		{
			name:        "successful create",
			requestBody: `{"name":"John","cpf":"12345678900","email":"john@example.com"}`,
			mockCreateFunc: func(ctx context.Context, dto dto.CreateCustomerDTO) (string, error) {
				return "id123", nil
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"id123","message":"Customer created successfully"}`,
		},
		{
			name:        "invalid json",
			requestBody: `{"name":"John",}`,
			mockCreateFunc: func(ctx context.Context, dto dto.CreateCustomerDTO) (string, error) {
				return "", nil
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid character '}' looking for beginning of object key string"}`,
		},
		{
			name:        "service create error",
			requestBody: `{"name":"John","cpf":"12345678900","email":"john@example.com"}`,
			mockCreateFunc: func(ctx context.Context, dto dto.CreateCustomerDTO) (string, error) {
				return "", errors.New("fail")
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"customer not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockCustomerService{
				createFunc: tt.mockCreateFunc,
			}

			handler := NewCustomerHandler(mockSvc)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/customer/register", handler.Create)

			req := httptest.NewRequest(http.MethodPost, "/customer/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != tt.expectedCode {
				t.Errorf("expected status %d but got %d", tt.expectedCode, resp.Code)
			}

			gotBody := strings.TrimSpace(resp.Body.String())
			expectedBody := strings.TrimSpace(tt.expectedBody)
			if gotBody != expectedBody {
				t.Errorf("expected body %s but got %s", expectedBody, gotBody)
			}
		})
	}
}

func TestCustomerHandler_Identify(t *testing.T) {
	tests := []struct {
		name             string
		paramCPF         string
		mockIdentifyFunc func(ctx context.Context, CPF string) (string, error)
		expectedCode     int
		expectedBody     string
	}{
		{
			name:     "successful identify",
			paramCPF: "12345678900",
			mockIdentifyFunc: func(ctx context.Context, CPF string) (string, error) {
				return "token123", nil
			},
			expectedCode: http.StatusOK,
			expectedBody: `"token123"`,
		},
		{
			name:     "identify error",
			paramCPF: "12345678900",
			mockIdentifyFunc: func(ctx context.Context, CPF string) (string, error) {
				return "", errors.New("not found")
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockCustomerService{
				identifyFunc: tt.mockIdentifyFunc,
			}

			handler := NewCustomerHandler(mockSvc)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.GET("/identify/:cpf", handler.Identify)

			req := httptest.NewRequest(http.MethodGet, "/identify/"+tt.paramCPF, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != tt.expectedCode {
				t.Errorf("expected status %d but got %d", tt.expectedCode, resp.Code)
			}

			gotBody := strings.TrimSpace(resp.Body.String())
			expectedBody := strings.TrimSpace(tt.expectedBody)
			if gotBody != expectedBody {
				t.Errorf("expected body %s but got %s", expectedBody, gotBody)
			}
		})
	}
}

func TestCustomerHandler_Anonymous(t *testing.T) {
	tests := []struct {
		name             string
		mockIdentifyFunc func(ctx context.Context, CPF string) (string, error)
		expectedCode     int
		expectedBody     string
	}{
		{
			name: "successful anonymous",
			mockIdentifyFunc: func(ctx context.Context, CPF string) (string, error) {
				if CPF != "" {
					t.Errorf("expected empty CPF but got %s", CPF)
				}
				return "tokenAnon", nil
			},
			expectedCode: http.StatusOK,
			expectedBody: `"tokenAnon"`,
		},
		{
			name: "anonymous error",
			mockIdentifyFunc: func(ctx context.Context, CPF string) (string, error) {
				return "", errors.New("not found")
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockCustomerService{
				identifyFunc: tt.mockIdentifyFunc,
			}

			handler := NewCustomerHandler(mockSvc)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.GET("/anonymous", handler.Anonymous)

			req := httptest.NewRequest(http.MethodGet, "/anonymous", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != tt.expectedCode {
				t.Errorf("expected status %d but got %d", tt.expectedCode, resp.Code)
			}

			gotBody := strings.TrimSpace(resp.Body.String())
			expectedBody := strings.TrimSpace(tt.expectedBody)
			if gotBody != expectedBody {
				t.Errorf("expected body %s but got %s", expectedBody, gotBody)
			}
		})
	}
}
