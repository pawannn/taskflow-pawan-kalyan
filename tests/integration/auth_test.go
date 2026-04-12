package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	auth "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/auth/jwt"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/logger"
	engine "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	authHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/auth"
	authService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/auth"
)

// mockUserRepo is an in-memory implementation of domainRepo.UserRepository.
type mockUserRepo struct {
	users map[string]*models.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*models.User)}
}

func (m *mockUserRepo) Create(_ context.Context, user *models.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*models.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepo) GetByID(_ context.Context, id string) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

// newAuthTestServer wires the full auth handler stack with an in-memory user repo.
func newAuthTestServer(t *testing.T) (*httptest.Server, *mockUserRepo) {
	t.Helper()

	cfg := &config.Config{
		AppName:    "taskflow-test",
		AppPort:    0,
		JWTSecret:  "test-secret-key-for-integration-tests",
		JWTExpiry:  24,
		BCryptCost: 4, // low cost for fast tests
	}

	log := logger.New()
	eng := engine.NewHttpEngine(cfg, log)

	tokenService := auth.NewTokenService(cfg.AppName, cfg.JWTSecret, cfg.JWTExpiry)
	userRepo := newMockUserRepo()
	svc := authService.NewAuthService(cfg.BCryptCost, userRepo, tokenService)

	handler := authHandler.NewAuthHandler(eng, svc)
	handler.AddRoutes()

	return httptest.NewServer(eng.Handler()), userRepo
}

// Test 1: POST /auth/register with blank required fields returns 400.
func TestRegister_MissingFields(t *testing.T) {
	srv, _ := newAuthTestServer(t)
	defer srv.Close()

	body := `{"name": "", "email": "", "password": ""}`
	resp, err := http.Post(srv.URL+"/auth/register", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result["error"] == nil {
		t.Error("expected error field in response body")
	}
}

// Test 2: POST /auth/register with valid data returns 201 and the created user.
func TestRegister_Success(t *testing.T) {
	srv, _ := newAuthTestServer(t)
	defer srv.Close()

	body := `{"name": "Jane Doe", "email": "jane@example.com", "password": "Test@Pass1!"}`
	resp, err := http.Post(srv.URL+"/auth/register", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in response")
	}
	if data["email"] != "jane@example.com" {
		t.Errorf("expected email jane@example.com, got %v", data["email"])
	}
	if data["name"] != "Jane Doe" {
		t.Errorf("expected name Jane Doe, got %v", data["name"])
	}
}

// Test 3: POST /auth/login with incorrect password returns 401.
func TestLogin_WrongPassword(t *testing.T) {
	srv, _ := newAuthTestServer(t)
	defer srv.Close()

	// Register a user first.
	regBody := `{"name": "Test User", "email": "test@example.com", "password": "Correct@Pass1!"}`
	regResp, err := http.Post(srv.URL+"/auth/register", "application/json", bytes.NewBufferString(regBody))
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}
	regResp.Body.Close()

	// Login with wrong password.
	loginBody := `{"email": "test@example.com", "password": "Wrong@Pass1!"}`
	resp, err := http.Post(srv.URL+"/auth/login", "application/json", bytes.NewBufferString(loginBody))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}
