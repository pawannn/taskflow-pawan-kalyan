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
	taskHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/task"
	middlewares "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/middlewares"
	authService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/auth"
	projectService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/project"
	taskService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/task"
)

type taskTestEnv struct {
	server      *httptest.Server
	userRepo    *mockUserRepo
	projectRepo *mockProjectRepo
	taskRepo    *mockTaskRepo
}

func newTaskTestEnv(t *testing.T) *taskTestEnv {
	t.Helper()

	cfg := &config.Config{
		Env:        "TEST",
		AppName:    "taskflow-test",
		JWTSecret:  "test-secret-for-task-tests",
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
	taskSvc := taskService.NewTaskService(taskRepo, projectRepo, userRepo)

	authHandler.NewAuthHandler(eng, authSvc).AddRoutes()
	projectHandler.NewProjectHandler(eng, projSvc, mw).AddRoutes()
	taskHandler.NewTaskHandler(eng, mw, taskSvc).AddRoutes()

	return &taskTestEnv{
		server:      httptest.NewServer(eng.Handler()),
		userRepo:    userRepo,
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
	}
}

// createProject is a helper that creates a project via the API and returns its ID.
func createProject(t *testing.T, serverURL, token, name string) string {
	t.Helper()

	body := fmt.Sprintf(`{"name":%q}`, name)
	req, _ := http.NewRequest(http.MethodPost, serverURL+"/projects",
		bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("createProject request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("createProject expected 201, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode createProject response: %v", err)
	}
	return result["data"].(map[string]interface{})["id"].(string)
}

// Test 1: POST /projects/:id/tasks without a token returns 401.
func TestCreateTask_Unauthenticated(t *testing.T) {
	env := newTaskTestEnv(t)
	defer env.server.Close()

	resp, err := http.Post(
		env.server.URL+"/projects/00000000-0000-0000-0000-000000000001/tasks",
		"application/json",
		bytes.NewBufferString(`{"title":"My Task"}`),
	)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

// Test 2: POST /projects/:id/tasks with a missing title returns 400 with a field error.
func TestCreateTask_MissingTitle(t *testing.T) {
	env := newTaskTestEnv(t)
	defer env.server.Close()

	token, _ := registerAndLogin(t, env.server.URL, "Eve", "eve@example.com", "Test@Pass1!")
	projectID := createProject(t, env.server.URL, token, "Sprint 1")

	req, _ := http.NewRequest(http.MethodPost,
		env.server.URL+"/projects/"+projectID+"/tasks",
		bytes.NewBufferString(`{"title":""}`))
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
	if !ok || fields["title"] == nil {
		t.Error("expected fields.title in error response")
	}
}

// Test 3: POST /projects/:id/tasks with a valid body returns 201 with the created task.
func TestCreateTask_Success(t *testing.T) {
	env := newTaskTestEnv(t)
	defer env.server.Close()

	token, _ := registerAndLogin(t, env.server.URL, "Frank", "frank@example.com", "Test@Pass1!")
	projectID := createProject(t, env.server.URL, token, "Release v2")

	body := `{"title":"Design homepage","priority":"high"}`
	req, _ := http.NewRequest(http.MethodPost,
		env.server.URL+"/projects/"+projectID+"/tasks",
		bytes.NewBufferString(body))
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
	if data["title"] != "Design homepage" {
		t.Errorf("expected title 'Design homepage', got %v", data["title"])
	}
	if data["status"] != "todo" {
		t.Errorf("expected default status 'todo', got %v", data["status"])
	}
	if data["id"] == nil || data["id"] == "" {
		t.Error("expected non-empty id in response")
	}
}

// Test 4: PATCH /tasks/:id with an invalid status value returns 400.
func TestUpdateTask_InvalidStatus(t *testing.T) {
	env := newTaskTestEnv(t)
	defer env.server.Close()

	token, _ := registerAndLogin(t, env.server.URL, "Grace", "grace@example.com", "Test@Pass1!")
	projectID := createProject(t, env.server.URL, token, "Backend")

	// Create a task.
	createReq, _ := http.NewRequest(http.MethodPost,
		env.server.URL+"/projects/"+projectID+"/tasks",
		bytes.NewBufferString(`{"title":"Setup DB"}`))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		t.Fatalf("create task failed: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	if err := json.NewDecoder(createResp.Body).Decode(&createResult); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	taskID := createResult["data"].(map[string]interface{})["id"].(string)

	// PATCH with a bad status.
	patchReq, _ := http.NewRequest(http.MethodPatch,
		env.server.URL+"/tasks/"+taskID,
		bytes.NewBufferString(`{"status":"invalid_status"}`))
	patchReq.Header.Set("Authorization", "Bearer "+token)
	patchReq.Header.Set("Content-Type", "application/json")

	patchResp, err := http.DefaultClient.Do(patchReq)
	if err != nil {
		t.Fatalf("patch request failed: %v", err)
	}
	defer patchResp.Body.Close()

	if patchResp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", patchResp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(patchResp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	fields, ok := result["fields"].(map[string]interface{})
	if !ok || fields["status"] == nil {
		t.Error("expected fields.status in error response")
	}
}

// Test 5: DELETE /tasks/:id by a user who is neither project owner nor task creator returns 403.
func TestDeleteTask_Forbidden(t *testing.T) {
	env := newTaskTestEnv(t)
	defer env.server.Close()

	// User A creates the project and task.
	tokenA, _ := registerAndLogin(t, env.server.URL, "Heidi", "heidi@example.com", "Test@Pass1!")
	projectID := createProject(t, env.server.URL, tokenA, "Heidi's Project")

	createReq, _ := http.NewRequest(http.MethodPost,
		env.server.URL+"/projects/"+projectID+"/tasks",
		bytes.NewBufferString(`{"title":"Secret Task"}`))
	createReq.Header.Set("Authorization", "Bearer "+tokenA)
	createReq.Header.Set("Content-Type", "application/json")
	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		t.Fatalf("create task failed: %v", err)
	}
	defer createResp.Body.Close()

	var createResult map[string]interface{}
	if err := json.NewDecoder(createResp.Body).Decode(&createResult); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	taskID := createResult["data"].(map[string]interface{})["id"].(string)

	// User B tries to delete User A's task.
	tokenB, _ := registerAndLogin(t, env.server.URL, "Ivan", "ivan@example.com", "Test@Pass1!")
	deleteReq, _ := http.NewRequest(http.MethodDelete,
		env.server.URL+"/tasks/"+taskID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+tokenB)

	deleteResp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", deleteResp.StatusCode)
	}
}
