package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/auth/jwt"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/config"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/infrastructure/logger"
	engine "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	authHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/auth"
	projectHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/project"
	middlewares "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/middlewares"
	authService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/auth"
	projectService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/project"
)

type projectTestEnv struct {
	server      *httptest.Server
	userRepo    *mockUserRepo
	projectRepo *mockProjectRepo
	taskRepo    *mockTaskRepo
}

func newProjectTestEnv(t *testing.T) *projectTestEnv {
	t.Helper()

	cfg := &config.Config{
		Env:        "TEST",
		AppName:    "taskflow-test",
		JWTSecret:  "test-secret-for-project-tests",
		JWTExpiry:  24,
		BCryptCost: 4,
	}

	log := logger.New(cfg.Env)
	eng := engine.NewHttpEngine(cfg, log)

	tokenService := auth.NewTokenService(cfg.AppName, cfg.JWTSecret, cfg.JWTExpiry)
	mw := middlewares.NewMiddlewareHandler(eng, *tokenService)

	userRepo := newMockUserRepo()
	projectRepo := newMockProjectRepo()
	taskRepo := newMockTaskRepo()

	authSvc := authService.NewAuthService(cfg.BCryptCost, userRepo, tokenService)
	projSvc := projectService.NewProjectService(projectRepo, taskRepo)

	authHandler.NewAuthHandler(eng, authSvc).AddRoutes()
	projectHandler.NewProjectHandler(eng, projSvc, mw).AddRoutes()

	return &projectTestEnv{
		server:      httptest.NewServer(eng.Handler()),
		userRepo:    userRepo,
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
	}
}

// registerAndLogin registers a user and returns the JWT token + user ID.
func registerAndLogin(t *testing.T, serverURL, name, email, password string) (token string, userID string) {
	t.Helper()

	regBody := fmt.Sprintf(`{"name":%q,"email":%q,"password":%q}`, name, email, password)
	regResp, err := http.Post(serverURL+"/auth/register", "application/json", bytes.NewBufferString(regBody))
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}
	defer regResp.Body.Close()
	if regResp.StatusCode != http.StatusCreated {
		t.Fatalf("register expected 201, got %d", regResp.StatusCode)
	}

	var regResult map[string]interface{}
	if err := json.NewDecoder(regResp.Body).Decode(&regResult); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	data := regResult["data"].(map[string]interface{})
	userID = data["id"].(string)

	loginBody := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)
	loginResp, err := http.Post(serverURL+"/auth/login", "application/json", bytes.NewBufferString(loginBody))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer loginResp.Body.Close()

	var loginResult map[string]interface{}
	if err := json.NewDecoder(loginResp.Body).Decode(&loginResult); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	loginData := loginResult["data"].(map[string]interface{})
	token = loginData["token"].(string)
	return token, userID
}

// Test 1: POST /projects without a token returns 401.
func TestCreateProject_Unauthenticated(t *testing.T) {
	env := newProjectTestEnv(t)
	defer env.server.Close()

	resp, err := http.Post(env.server.URL+"/projects", "application/json",
		bytes.NewBufferString(`{"name":"My Project"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

// Test 2: POST /projects with an empty name returns 400 with a validation field error.
func TestCreateProject_MissingName(t *testing.T) {
	env := newProjectTestEnv(t)
	defer env.server.Close()

	token, _ := registerAndLogin(t, env.server.URL, "Alice", "alice@example.com", "Test@Pass1!")

	req, _ := http.NewRequest(http.MethodPost, env.server.URL+"/projects",
		bytes.NewBufferString(`{"name":""}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	fields, ok := result["fields"].(map[string]interface{})
	if !ok || fields["name"] == nil {
		t.Error("expected fields.name in error response")
	}
}

// Test 3: POST /projects with a valid body returns 201 and the created project.
func TestCreateProject_Success(t *testing.T) {
	env := newProjectTestEnv(t)
	defer env.server.Close()

	token, userID := registerAndLogin(t, env.server.URL, "Bob", "bob@example.com", "Test@Pass1!")

	req, _ := http.NewRequest(http.MethodPost, env.server.URL+"/projects",
		bytes.NewBufferString(`{"name":"Website Redesign","description":"Q2 initiative"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data := result["data"].(map[string]interface{})
	if data["name"] != "Website Redesign" {
		t.Errorf("expected name 'Website Redesign', got %v", data["name"])
	}
	if data["owner_id"] != userID {
		t.Errorf("expected owner_id %s, got %v", userID, data["owner_id"])
	}
	if data["id"] == nil || data["id"] == "" {
		t.Error("expected non-empty id in response")
	}
}

// Test 4: GET /projects returns projects owned by the authenticated user.
func TestGetProjects_ReturnsOwnedProjects(t *testing.T) {
	env := newProjectTestEnv(t)
	defer env.server.Close()

	token, _ := registerAndLogin(t, env.server.URL, "Carol", "carol@example.com", "Test@Pass1!")

	// Create two projects.
	for _, name := range []string{"Alpha", "Beta"} {
		body := fmt.Sprintf(`{"name":%q}`, name)
		req, _ := http.NewRequest(http.MethodPost, env.server.URL+"/projects",
			bytes.NewBufferString(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		if _, err := http.DefaultClient.Do(req); err != nil {
			t.Fatalf("create project %q failed: %v", name, err)
		}
	}

	// List projects.
	req, _ := http.NewRequest(http.MethodGet, env.server.URL+"/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data := result["data"].(map[string]interface{})
	projects := data["projects"].([]interface{})
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
}

// Test 5: PATCH /projects/:id by a user who is not the owner returns 403.
func TestUpdateProject_Forbidden(t *testing.T) {
	env := newProjectTestEnv(t)
	defer env.server.Close()

	// User A creates a project.
	tokenA, _ := registerAndLogin(t, env.server.URL, "Dave", "dave@example.com", "Test@Pass1!")
	createReq, _ := http.NewRequest(http.MethodPost, env.server.URL+"/projects",
		bytes.NewBufferString(`{"name":"Dave's Project"}`))
	createReq.Header.Set("Authorization", "Bearer "+tokenA)
	createReq.Header.Set("Content-Type", "application/json")
	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		t.Fatalf("create project failed: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	if err := json.NewDecoder(createResp.Body).Decode(&createResult); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	projectID := createResult["data"].(map[string]interface{})["id"].(string)

	// User B tries to PATCH User A's project.
	tokenB, _ := registerAndLogin(t, env.server.URL, "Eve", "eve@example.com", "Test@Pass1!")
	patchReq, _ := http.NewRequest(http.MethodPatch,
		env.server.URL+"/projects/"+projectID,
		bytes.NewBufferString(`{"name":"Hijacked"}`))
	patchReq.Header.Set("Authorization", "Bearer "+tokenB)
	patchReq.Header.Set("Content-Type", "application/json")

	patchResp, err := http.DefaultClient.Do(patchReq)
	if err != nil {
		t.Fatalf("patch request failed: %v", err)
	}
	defer patchResp.Body.Close()

	if patchResp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", patchResp.StatusCode)
	}
}
