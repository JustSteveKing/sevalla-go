package sevalla

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestApplicationsService_List(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create client with test server URL
	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	// Expected applications
	want := []*Application{
		{
			ID:     "app-1",
			Name:   "test-app-1",
			State:  StateRunning,
			Region: RegionUSCentral,
			Plan:   PlanHobby,
		},
		{
			ID:     "app-2",
			Name:   "test-app-2",
			State:  StateDeploying,
			Region: RegionEuropeWest,
			Plan:   PlanPro,
		},
	}

	// Setup endpoint handler
	mux.HandleFunc("/applications", func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check authentication header
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got %s", auth)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	// Call the method
	ctx := context.Background()
	apps, _, err := client.Applications.List(ctx, nil)

	// Check results
	if err != nil {
		t.Fatalf("Applications.List returned error: %v", err)
	}

	if !reflect.DeepEqual(apps, want) {
		t.Errorf("Applications.List returned %+v, want %+v", apps, want)
	}
}

func TestApplicationsService_Create(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	// Request to create
	createReq := &CreateApplicationRequest{
		Name:          "new-app",
		RepositoryURL: "https://github.com/user/repo",
		Branch:        "main",
		Region:        RegionUSCentral,
		Plan:          PlanStarter,
		Replicas:      2,
	}

	// Expected response
	want := &Application{
		ID:            "app-new",
		Name:          "new-app",
		State:         StatePending,
		RepositoryURL: "https://github.com/user/repo",
		Branch:        "main",
		Region:        RegionUSCentral,
		Plan:          PlanStarter,
		Replicas:      2,
		CreatedAt:     time.Now(),
	}

	// Setup endpoint handler
	mux.HandleFunc("/applications", func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check request body
		var got CreateApplicationRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if !reflect.DeepEqual(&got, createReq) {
			t.Errorf("Request body = %+v, want %+v", got, createReq)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	// Call the method
	ctx := context.Background()
	app, resp, err := client.Applications.Create(ctx, createReq)

	// Check results
	if err != nil {
		t.Fatalf("Applications.Create returned error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", resp.StatusCode)
	}

	if app.ID != want.ID || app.Name != want.Name {
		t.Errorf("Applications.Create returned %+v, want %+v", app, want)
	}
}

func TestApplicationsService_Deploy(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	appID := "app-123"

	// Expected deployment
	want := &Deployment{
		ID:            "deploy-456",
		ApplicationID: appID,
		State:         StatusQueued,
		Branch:        "main",
		StartedAt:     time.Now(),
	}

	// Setup endpoint handler
	mux.HandleFunc("/applications/"+appID+"/deployments", func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	// Call the method
	ctx := context.Background()
	deployment, _, err := client.Applications.Deploy(ctx, appID)

	// Check results
	if err != nil {
		t.Fatalf("Applications.Deploy returned error: %v", err)
	}

	if deployment.ID != want.ID || deployment.ApplicationID != want.ApplicationID {
		t.Errorf("Applications.Deploy returned %+v, want %+v", deployment, want)
	}
}

func TestApplicationsService_Update(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	updateReq := &UpdateApplicationRequest{
		Name:     String("updated-app"),
		Replicas: Int(3),
	}

	want := &Application{
		ID:       "app-123",
		Name:     "updated-app",
		Replicas: 3,
	}

	mux.HandleFunc("/applications/app-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	app, _, err := client.Applications.Update(ctx, "app-123", updateReq)

	if err != nil {
		t.Fatalf("Applications.Update returned error: %v", err)
	}

	if app.Name != want.Name {
		t.Errorf("Expected app name %s, got %s", want.Name, app.Name)
	}
}

func TestApplicationsService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/applications/app-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Applications.Delete(ctx, "app-123")

	if err != nil {
		t.Fatalf("Applications.Delete returned error: %v", err)
	}
}

func TestApplicationsService_Scale(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	plan := PlanPro
	scaleReq := &ScaleApplicationRequest{
		Replicas: 5,
		Plan:     &plan,
	}

	want := &Application{
		ID:       "app-123",
		Replicas: 5,
		Plan:     PlanPro,
	}

	mux.HandleFunc("/applications/app-123/scale", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	app, _, err := client.Applications.Scale(ctx, "app-123", scaleReq)

	if err != nil {
		t.Fatalf("Applications.Scale returned error: %v", err)
	}

	if app.Replicas != want.Replicas {
		t.Errorf("Expected replicas %d, got %d", want.Replicas, app.Replicas)
	}
}

func TestApplicationsService_Restart(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/applications/app-123/restart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	_, err := client.Applications.Restart(ctx, "app-123")

	if err != nil {
		t.Fatalf("Applications.Restart returned error: %v", err)
	}
}

func TestApplicationsService_Stop(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/applications/app-123/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	_, err := client.Applications.Stop(ctx, "app-123")

	if err != nil {
		t.Fatalf("Applications.Stop returned error: %v", err)
	}
}

func TestApplicationsService_Start(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/applications/app-123/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	_, err := client.Applications.Start(ctx, "app-123")

	if err != nil {
		t.Fatalf("Applications.Start returned error: %v", err)
	}
}

func TestApplicationsService_GetLogs(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	wantLogs := "Application logs here\nLine 2\nLine 3"

	mux.HandleFunc("/applications/app-123/logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check query parameter
		if r.URL.Query().Get("lines") != "100" {
			t.Errorf("Expected lines=100 query parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"logs": wantLogs}); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	logs, _, err := client.Applications.GetLogs(ctx, "app-123", 100)

	if err != nil {
		t.Fatalf("Applications.GetLogs returned error: %v", err)
	}

	if logs != wantLogs {
		t.Errorf("Expected logs %s, got %s", wantLogs, logs)
	}
}

func TestApplicationsService_ListDeployments(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := []*Deployment{
		{
			ID:            "deploy-1",
			ApplicationID: "app-123",
			State:         StatusSuccess,
		},
		{
			ID:            "deploy-2",
			ApplicationID: "app-123",
			State:         StatusFailed,
		},
	}

	mux.HandleFunc("/applications/app-123/deployments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	deployments, _, err := client.Applications.ListDeployments(ctx, "app-123", nil)

	if err != nil {
		t.Fatalf("Applications.ListDeployments returned error: %v", err)
	}

	if len(deployments) != len(want) {
		t.Errorf("Expected %d deployments, got %d", len(want), len(deployments))
	}
}

func TestApplicationsService_GetDeployment(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Deployment{
		ID:            "deploy-1",
		ApplicationID: "app-123",
		State:         StatusSuccess,
		Branch:        "main",
	}

	mux.HandleFunc("/applications/app-123/deployments/deploy-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	deployment, _, err := client.Applications.GetDeployment(ctx, "app-123", "deploy-1")

	if err != nil {
		t.Fatalf("Applications.GetDeployment returned error: %v", err)
	}

	if deployment.ID != want.ID {
		t.Errorf("Expected deployment ID %s, got %s", want.ID, deployment.ID)
	}
}

func TestApplicationsService_CancelDeployment(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/applications/app-123/deployments/deploy-1/cancel", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	_, err := client.Applications.CancelDeployment(ctx, "app-123", "deploy-1")

	if err != nil {
		t.Fatalf("Applications.CancelDeployment returned error: %v", err)
	}
}

func TestApplicationsService_AddCustomDomain(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/applications/app-123/domains", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var req AddDomainRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		if req.Domain != "example.com" {
			t.Errorf("Expected domain example.com, got %s", req.Domain)
		}

		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	_, err := client.Applications.AddCustomDomain(ctx, "app-123", "example.com")

	if err != nil {
		t.Fatalf("Applications.AddCustomDomain returned error: %v", err)
	}
}

func TestApplicationsService_RemoveCustomDomain(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/applications/app-123/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Applications.RemoveCustomDomain(ctx, "app-123", "example.com")

	if err != nil {
		t.Fatalf("Applications.RemoveCustomDomain returned error: %v", err)
	}
}

func TestApplicationsService_UpdateCDNSettings(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/applications/app-123/cdn", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		var req CDNSettingsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		if !req.Enabled {
			t.Error("Expected CDN to be enabled")
		}

		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	_, err := client.Applications.UpdateCDNSettings(ctx, "app-123", true)

	if err != nil {
		t.Fatalf("Applications.UpdateCDNSettings returned error: %v", err)
	}
}

func TestApplicationsService_GetUsage(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Usage{
		ApplicationID: "app-123",
		CPUUsage:      75.5,
		MemoryUsage:   2048,
		StorageUsage:  10000,
	}

	mux.HandleFunc("/applications/app-123/usage", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		if r.URL.Query().Get("period") != "7d" {
			t.Errorf("Expected period=7d query parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	usage, _, err := client.Applications.GetUsage(ctx, "app-123", "7d")

	if err != nil {
		t.Fatalf("Applications.GetUsage returned error: %v", err)
	}

	if usage.ApplicationID != want.ApplicationID {
		t.Errorf("Expected application ID %s, got %s", want.ApplicationID, usage.ApplicationID)
	}
}

func TestApplicationsService_SetEnvironmentVariables(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	vars := map[string]string{
		"DATABASE_URL": "postgresql://localhost/db",
		"API_KEY":      "secret123",
	}

	mux.HandleFunc("/applications/app-123/env", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		var reqVars map[string]string
		if err := json.NewDecoder(r.Body).Decode(&reqVars); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		if reqVars["DATABASE_URL"] != vars["DATABASE_URL"] {
			t.Error("Environment variables mismatch")
		}

		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	_, err := client.Applications.SetEnvironmentVariables(ctx, "app-123", vars)

	if err != nil {
		t.Fatalf("Applications.SetEnvironmentVariables returned error: %v", err)
	}
}

func TestApplicationsService_GetEnvironmentVariables(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := map[string]string{
		"DATABASE_URL": "postgresql://localhost/db",
		"API_KEY":      "secret123",
	}

	mux.HandleFunc("/applications/app-123/env", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	vars, _, err := client.Applications.GetEnvironmentVariables(ctx, "app-123")

	if err != nil {
		t.Fatalf("Applications.GetEnvironmentVariables returned error: %v", err)
	}

	if vars["DATABASE_URL"] != want["DATABASE_URL"] {
		t.Errorf("Expected DATABASE_URL %s, got %s", want["DATABASE_URL"], vars["DATABASE_URL"])
	}
}

func TestApplicationsService_Rollback(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Deployment{
		ID:            "deploy-new",
		ApplicationID: "app-123",
		State:         StatusQueued,
	}

	mux.HandleFunc("/applications/app-123/rollback/deploy-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	deployment, _, err := client.Applications.Rollback(ctx, "app-123", "deploy-1")

	if err != nil {
		t.Fatalf("Applications.Rollback returned error: %v", err)
	}

	if deployment.ID != want.ID {
		t.Errorf("Expected deployment ID %s, got %s", want.ID, deployment.ID)
	}
}

func TestErrorResponse(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	// Setup endpoint to return error
	mux.HandleFunc("/applications/not-found", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Application not found",
			"code":    "RESOURCE_NOT_FOUND",
		}); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	// Call the method
	ctx := context.Background()
	_, _, err := client.Applications.Get(ctx, "not-found")

	// Check error
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error type
	if !IsNotFound(err) {
		t.Errorf("Expected IsNotFound to return true, got false")
	}

	// Check error message
	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("Expected ErrorResponse type, got %T", err)
	}

	if errResp.Message != "Application not found" {
		t.Errorf("Expected error message 'Application not found', got %s", errResp.Message)
	}

	if errResp.Code != "RESOURCE_NOT_FOUND" {
		t.Errorf("Expected error code 'RESOURCE_NOT_FOUND', got %s", errResp.Code)
	}
}

func TestPaginationResponse(t *testing.T) {
	// Setup test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	// Setup endpoint with Link header for pagination
	mux.HandleFunc("/applications", func(w http.ResponseWriter, r *http.Request) {
		// Set Link header with pagination
		linkHeader := `<` + server.URL + `/applications?page=2>; rel="next", ` +
			`<` + server.URL + `/applications?page=1>; rel="prev", ` +
			`<` + server.URL + `/applications?page=1>; rel="first", ` +
			`<` + server.URL + `/applications?page=5>; rel="last"`

		w.Header().Set("Link", linkHeader)
		w.Header().Set("Content-Type", "application/json")

		// Return empty array
		if err := json.NewEncoder(w).Encode([]*Application{}); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	// Call the method
	ctx := context.Background()
	_, resp, err := client.Applications.List(ctx, nil)

	// Check results
	if err != nil {
		t.Fatalf("Applications.List returned error: %v", err)
	}

	// Check pagination values
	if resp.NextPage != 2 {
		t.Errorf("Expected NextPage = 2, got %d", resp.NextPage)
	}

	if resp.PrevPage != 1 {
		t.Errorf("Expected PrevPage = 1, got %d", resp.PrevPage)
	}

	if resp.FirstPage != 1 {
		t.Errorf("Expected FirstPage = 1, got %d", resp.FirstPage)
	}

	if resp.LastPage != 5 {
		t.Errorf("Expected LastPage = 5, got %d", resp.LastPage)
	}
}

func TestClientOptions(t *testing.T) {
	// Test with custom HTTP client
	customHTTPClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	client := NewClient(
		WithAPIKey("custom-key"),
		WithHTTPClient(customHTTPClient),
		WithBaseURL("https://custom.api.url"),
		WithUserAgent("custom-agent/1.0"),
	)

	// Check API key
	if client.apiKey != "custom-key" {
		t.Errorf("Expected API key 'custom-key', got %s", client.apiKey)
	}

	// Check HTTP client
	if client.client != customHTTPClient {
		t.Error("Expected custom HTTP client to be set")
	}

	// Check base URL
	if client.baseURL.String() != "https://custom.api.url" {
		t.Errorf("Expected base URL 'https://custom.api.url', got %s", client.baseURL.String())
	}

	// Check user agent
	if client.userAgent != "custom-agent/1.0" {
		t.Errorf("Expected user agent 'custom-agent/1.0', got %s", client.userAgent)
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test Bool helper
	boolVal := true
	boolPtr := Bool(boolVal)
	if *boolPtr != boolVal {
		t.Errorf("Bool helper returned %v, want %v", *boolPtr, boolVal)
	}

	// Test Int helper
	intVal := 42
	intPtr := Int(intVal)
	if *intPtr != intVal {
		t.Errorf("Int helper returned %v, want %v", *intPtr, intVal)
	}

	// Test String helper
	strVal := "test"
	strPtr := String(strVal)
	if *strPtr != strVal {
		t.Errorf("String helper returned %v, want %v", *strPtr, strVal)
	}
}

// Database Service Tests

func TestDatabasesService_List(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := []*Database{
		{
			ID:      "db-1",
			Name:    "test-db-1",
			Type:    EnginePostgreSQL,
			Version: "14",
			Region:  RegionUSCentral,
		},
	}

	mux.HandleFunc("/databases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	databases, _, err := client.Databases.List(ctx, nil)

	if err != nil {
		t.Fatalf("Databases.List returned error: %v", err)
	}

	if len(databases) != len(want) {
		t.Errorf("Expected %d databases, got %d", len(want), len(databases))
	}
}

func TestDatabasesService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Database{
		ID:      "db-1",
		Name:    "test-db",
		Type:    EngineMySQL,
		Version: "8.0",
		Region:  RegionUSCentral,
	}

	mux.HandleFunc("/databases/db-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	db, _, err := client.Databases.Get(ctx, "db-1")

	if err != nil {
		t.Fatalf("Databases.Get returned error: %v", err)
	}

	if db.ID != want.ID {
		t.Errorf("Expected database ID %s, got %s", want.ID, db.ID)
	}
}

func TestDatabasesService_Create(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	createReq := &CreateDatabaseRequest{
		Name:    "new-db",
		Type:    EnginePostgreSQL,
		Version: "14",
		Region:  RegionUSCentral,
	}

	want := &Database{
		ID:      "db-new",
		Name:    "new-db",
		Type:    EnginePostgreSQL,
		Version: "14",
		Region:  RegionUSCentral,
	}

	mux.HandleFunc("/databases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	db, _, err := client.Databases.Create(ctx, createReq)

	if err != nil {
		t.Fatalf("Databases.Create returned error: %v", err)
	}

	if db.ID != want.ID {
		t.Errorf("Expected database ID %s, got %s", want.ID, db.ID)
	}
}

func TestDatabasesService_Update(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	updateReq := &UpdateDatabaseRequest{
		Name: String("updated-db"),
	}

	want := &Database{
		ID:   "db-1",
		Name: "updated-db",
	}

	mux.HandleFunc("/databases/db-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	db, _, err := client.Databases.Update(ctx, "db-1", updateReq)

	if err != nil {
		t.Fatalf("Databases.Update returned error: %v", err)
	}

	if db.Name != want.Name {
		t.Errorf("Expected database name %s, got %s", want.Name, db.Name)
	}
}

func TestDatabasesService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/databases/db-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Databases.Delete(ctx, "db-1")

	if err != nil {
		t.Fatalf("Databases.Delete returned error: %v", err)
	}
}

func TestDatabasesService_GetCredentials(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Database{
		ID:          "db-1",
		InternalURL: "postgresql://db.example.com:5432",
		Username:    "dbuser",
		Password:    "secret",
	}

	mux.HandleFunc("/databases/db-1/credentials", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	db, _, err := client.Databases.GetCredentials(ctx, "db-1")

	if err != nil {
		t.Fatalf("Databases.GetCredentials returned error: %v", err)
	}

	if db.Username != want.Username {
		t.Errorf("Expected username %s, got %s", want.Username, db.Username)
	}
}

func TestDatabasesService_ResetPassword(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Database{
		ID:       "db-1",
		Password: "new-secret",
	}

	mux.HandleFunc("/databases/db-1/reset-password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	db, _, err := client.Databases.ResetPassword(ctx, "db-1")

	if err != nil {
		t.Fatalf("Databases.ResetPassword returned error: %v", err)
	}

	if db.Password != want.Password {
		t.Errorf("Expected password %s, got %s", want.Password, db.Password)
	}
}

func TestDatabasesService_ListBackups(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := []*Backup{
		{
			ID:         "backup-1",
			DatabaseID: "db-1",
			Type:       "manual",
			Status:     "completed",
		},
	}

	mux.HandleFunc("/databases/db-1/backups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	backups, _, err := client.Databases.ListBackups(ctx, "db-1", nil)

	if err != nil {
		t.Fatalf("Databases.ListBackups returned error: %v", err)
	}

	if len(backups) != len(want) {
		t.Errorf("Expected %d backups, got %d", len(want), len(backups))
	}
}

func TestDatabasesService_CreateBackup(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	backupReq := &CreateBackupRequest{
		Type:        "manual",
		Description: "Test backup",
	}

	want := &Backup{
		ID:         "backup-new",
		DatabaseID: "db-1",
		Type:       "manual",
		Status:     "pending",
	}

	mux.HandleFunc("/databases/db-1/backups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	backup, _, err := client.Databases.CreateBackup(ctx, "db-1", backupReq)

	if err != nil {
		t.Fatalf("Databases.CreateBackup returned error: %v", err)
	}

	if backup.ID != want.ID {
		t.Errorf("Expected backup ID %s, got %s", want.ID, backup.ID)
	}
}

func TestDatabasesService_GetBackup(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Backup{
		ID:         "backup-1",
		DatabaseID: "db-1",
		Type:       "manual",
		Status:     "completed",
	}

	mux.HandleFunc("/databases/db-1/backups/backup-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	backup, _, err := client.Databases.GetBackup(ctx, "db-1", "backup-1")

	if err != nil {
		t.Fatalf("Databases.GetBackup returned error: %v", err)
	}

	if backup.ID != want.ID {
		t.Errorf("Expected backup ID %s, got %s", want.ID, backup.ID)
	}
}

func TestDatabasesService_DeleteBackup(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/databases/db-1/backups/backup-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Databases.DeleteBackup(ctx, "db-1", "backup-1")

	if err != nil {
		t.Fatalf("Databases.DeleteBackup returned error: %v", err)
	}
}

func TestDatabasesService_RestoreFromBackup(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	restoreReq := &RestoreBackupRequest{
		BackupID: "backup-1",
	}

	mux.HandleFunc("/databases/db-1/restore", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	_, err := client.Databases.RestoreFromBackup(ctx, "db-1", restoreReq)

	if err != nil {
		t.Fatalf("Databases.RestoreFromBackup returned error: %v", err)
	}
}

func TestDatabasesService_GetUsage(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Usage{
		DatabaseID:   "db-1",
		CPUUsage:     50.5,
		MemoryUsage:  1024,
		StorageUsage: 5000,
	}

	mux.HandleFunc("/databases/db-1/usage", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	usage, _, err := client.Databases.GetUsage(ctx, "db-1", "7d")

	if err != nil {
		t.Fatalf("Databases.GetUsage returned error: %v", err)
	}

	if usage.CPUUsage != want.CPUUsage {
		t.Errorf("Expected CPU %f, got %f", want.CPUUsage, usage.CPUUsage)
	}
}

func TestDatabasesService_EnablePublicAccess(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Database{
		ID:        "db-1",
		PublicURL: "postgresql://public.example.com:5432",
	}

	mux.HandleFunc("/databases/db-1/public-access", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	db, _, err := client.Databases.EnablePublicAccess(ctx, "db-1")

	if err != nil {
		t.Fatalf("Databases.EnablePublicAccess returned error: %v", err)
	}

	if db.PublicURL == "" {
		t.Error("Expected public URL to be set")
	}
}

func TestDatabasesService_DisablePublicAccess(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Database{
		ID:        "db-1",
		PublicURL: "",
	}

	mux.HandleFunc("/databases/db-1/public-access", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	db, _, err := client.Databases.DisablePublicAccess(ctx, "db-1")

	if err != nil {
		t.Fatalf("Databases.DisablePublicAccess returned error: %v", err)
	}

	if db.PublicURL != "" {
		t.Error("Expected public URL to be empty")
	}
}

// Static Sites Service Tests

func TestStaticSitesService_List(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := []*StaticSite{
		{
			ID:     "site-1",
			Name:   "test-site",
			State:  StateRunning,
			Region: RegionUSCentral,
		},
	}

	mux.HandleFunc("/static-sites", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	sites, _, err := client.StaticSites.List(ctx, nil)

	if err != nil {
		t.Fatalf("StaticSites.List returned error: %v", err)
	}

	if len(sites) != len(want) {
		t.Errorf("Expected %d sites, got %d", len(want), len(sites))
	}
}

func TestStaticSitesService_ListWithCompanyID(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := []*StaticSite{
		{
			ID:     "site-1",
			Name:   "test-site",
			State:  StateRunning,
			Region: RegionUSCentral,
		},
	}

	mux.HandleFunc("/static-sites", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check that company_id query parameter is present
		companyID := r.URL.Query().Get("company_id")
		if companyID != "company-123" {
			t.Errorf("Expected company_id=company-123 query parameter, got %s", companyID)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	opts := &ListOptions{
		CompanyID: "company-123",
	}
	sites, _, err := client.StaticSites.List(ctx, opts)

	if err != nil {
		t.Fatalf("StaticSites.List returned error: %v", err)
	}

	if len(sites) != len(want) {
		t.Errorf("Expected %d sites, got %d", len(want), len(sites))
	}
}

func TestStaticSitesService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &StaticSite{
		ID:     "site-1",
		Name:   "test-site",
		State:  StateRunning,
		Region: RegionUSCentral,
	}

	mux.HandleFunc("/static-sites/site-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	site, _, err := client.StaticSites.Get(ctx, "site-1")

	if err != nil {
		t.Fatalf("StaticSites.Get returned error: %v", err)
	}

	if site.ID != want.ID {
		t.Errorf("Expected site ID %s, got %s", want.ID, site.ID)
	}
}

func TestStaticSitesService_Create(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	createReq := &CreateStaticSiteRequest{
		Name:          "new-site",
		RepositoryURL: "https://github.com/user/repo",
		Branch:        "main",
		Region:        RegionUSCentral,
	}

	want := &StaticSite{
		ID:            "site-new",
		Name:          "new-site",
		RepositoryURL: "https://github.com/user/repo",
		Branch:        "main",
		State:         StatePending,
		Region:        RegionUSCentral,
	}

	mux.HandleFunc("/static-sites", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	site, _, err := client.StaticSites.Create(ctx, createReq)

	if err != nil {
		t.Fatalf("StaticSites.Create returned error: %v", err)
	}

	if site.ID != want.ID {
		t.Errorf("Expected site ID %s, got %s", want.ID, site.ID)
	}
}

func TestStaticSitesService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/static-sites/site-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.StaticSites.Delete(ctx, "site-1")

	if err != nil {
		t.Fatalf("StaticSites.Delete returned error: %v", err)
	}
}

func TestStaticSitesService_Deploy(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Deployment{
		ID:    "deploy-1",
		State: StatusQueued,
	}

	mux.HandleFunc("/static-sites/site-1/deployments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	deployment, _, err := client.StaticSites.Deploy(ctx, "site-1")

	if err != nil {
		t.Fatalf("StaticSites.Deploy returned error: %v", err)
	}

	if deployment.ID != want.ID {
		t.Errorf("Expected deployment ID %s, got %s", want.ID, deployment.ID)
	}
}

// Deployments Service Tests

func TestDeploymentsService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Deployment{
		ID:    "deploy-1",
		State: StatusSuccess,
	}

	mux.HandleFunc("/deployments/deploy-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	deployment, _, err := client.Deployments.Get(ctx, "deploy-1")

	if err != nil {
		t.Fatalf("Deployments.Get returned error: %v", err)
	}

	if deployment.ID != want.ID {
		t.Errorf("Expected deployment ID %s, got %s", want.ID, deployment.ID)
	}
}

func TestDeploymentsService_List(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := []*Deployment{
		{
			ID:    "deploy-1",
			State: StatusSuccess,
		},
	}

	mux.HandleFunc("/deployments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	deployments, _, err := client.Deployments.List(ctx, nil)

	if err != nil {
		t.Fatalf("Deployments.List returned error: %v", err)
	}

	if len(deployments) != len(want) {
		t.Errorf("Expected %d deployments, got %d", len(want), len(deployments))
	}
}

func TestDeploymentsService_GetLogs(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	wantLogs := "deployment logs here"

	mux.HandleFunc("/deployments/deploy-1/logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"logs": wantLogs}); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	logs, _, err := client.Deployments.GetLogs(ctx, "deploy-1")

	if err != nil {
		t.Fatalf("Deployments.GetLogs returned error: %v", err)
	}

	if logs != wantLogs {
		t.Errorf("Expected logs %s, got %s", wantLogs, logs)
	}
}

func TestDeploymentsService_Cancel(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/deployments/deploy-1/cancel", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	_, err := client.Deployments.Cancel(ctx, "deploy-1")

	if err != nil {
		t.Fatalf("Deployments.Cancel returned error: %v", err)
	}
}

// Pipelines Service Tests

func TestPipelinesService_List(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := []*Pipeline{
		{
			ID:      "pipeline-1",
			Name:    "test-pipeline",
			Enabled: true,
			Branch:  "main",
		},
	}

	mux.HandleFunc("/pipelines", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	pipelines, _, err := client.Pipelines.List(ctx, nil)

	if err != nil {
		t.Fatalf("Pipelines.List returned error: %v", err)
	}

	if len(pipelines) != len(want) {
		t.Errorf("Expected %d pipelines, got %d", len(want), len(pipelines))
	}
}

func TestPipelinesService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &Pipeline{
		ID:      "pipeline-1",
		Name:    "test-pipeline",
		Enabled: true,
		Branch:  "main",
	}

	mux.HandleFunc("/pipelines/pipeline-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	pipeline, _, err := client.Pipelines.Get(ctx, "pipeline-1")

	if err != nil {
		t.Fatalf("Pipelines.Get returned error: %v", err)
	}

	if pipeline.ID != want.ID {
		t.Errorf("Expected pipeline ID %s, got %s", want.ID, pipeline.ID)
	}
}

func TestPipelinesService_Create(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	createReq := &CreatePipelineRequest{
		Name:    "new-pipeline",
		Enabled: true,
		Trigger: "push",
		Branch:  "main",
		Steps:   []PipelineStep{},
	}

	want := &Pipeline{
		ID:      "pipeline-new",
		Name:    "new-pipeline",
		Enabled: true,
		Branch:  "main",
	}

	mux.HandleFunc("/pipelines", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	pipeline, _, err := client.Pipelines.Create(ctx, createReq)

	if err != nil {
		t.Fatalf("Pipelines.Create returned error: %v", err)
	}

	if pipeline.ID != want.ID {
		t.Errorf("Expected pipeline ID %s, got %s", want.ID, pipeline.ID)
	}
}

func TestPipelinesService_Update(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	updateReq := &UpdatePipelineRequest{
		Name: String("updated-pipeline"),
	}

	want := &Pipeline{
		ID:   "pipeline-1",
		Name: "updated-pipeline",
	}

	mux.HandleFunc("/pipelines/pipeline-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	pipeline, _, err := client.Pipelines.Update(ctx, "pipeline-1", updateReq)

	if err != nil {
		t.Fatalf("Pipelines.Update returned error: %v", err)
	}

	if pipeline.Name != want.Name {
		t.Errorf("Expected pipeline name %s, got %s", want.Name, pipeline.Name)
	}
}

func TestPipelinesService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/pipelines/pipeline-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Pipelines.Delete(ctx, "pipeline-1")

	if err != nil {
		t.Fatalf("Pipelines.Delete returned error: %v", err)
	}
}

func TestPipelinesService_Run(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &PipelineRun{
		ID:         "run-1",
		PipelineID: "pipeline-1",
		State:      "pending",
	}

	mux.HandleFunc("/pipelines/pipeline-1/runs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	run, _, err := client.Pipelines.Run(ctx, "pipeline-1")

	if err != nil {
		t.Fatalf("Pipelines.Run returned error: %v", err)
	}

	if run.ID != want.ID {
		t.Errorf("Expected run ID %s, got %s", want.ID, run.ID)
	}
}

func TestPipelinesService_ListRuns(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := []*PipelineRun{
		{
			ID:         "run-1",
			PipelineID: "pipeline-1",
			State:      "completed",
		},
	}

	mux.HandleFunc("/pipelines/pipeline-1/runs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	runs, _, err := client.Pipelines.ListRuns(ctx, "pipeline-1", nil)

	if err != nil {
		t.Fatalf("Pipelines.ListRuns returned error: %v", err)
	}

	if len(runs) != len(want) {
		t.Errorf("Expected %d runs, got %d", len(want), len(runs))
	}
}

func TestPipelinesService_GetRun(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &PipelineRun{
		ID:         "run-1",
		PipelineID: "pipeline-1",
		State:      "completed",
	}

	mux.HandleFunc("/pipelines/pipeline-1/runs/run-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	run, _, err := client.Pipelines.GetRun(ctx, "pipeline-1", "run-1")

	if err != nil {
		t.Fatalf("Pipelines.GetRun returned error: %v", err)
	}

	if run.ID != want.ID {
		t.Errorf("Expected run ID %s, got %s", want.ID, run.ID)
	}
}

func TestPipelinesService_CancelRun(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/pipelines/pipeline-1/runs/run-1/cancel", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	_, err := client.Pipelines.CancelRun(ctx, "pipeline-1", "run-1")

	if err != nil {
		t.Fatalf("Pipelines.CancelRun returned error: %v", err)
	}
}

func TestPipelinesService_GetRunLogs(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	wantLogs := "pipeline run logs here"

	mux.HandleFunc("/pipelines/pipeline-1/runs/run-1/logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"logs": wantLogs}); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	logs, _, err := client.Pipelines.GetRunLogs(ctx, "pipeline-1", "run-1")

	if err != nil {
		t.Fatalf("Pipelines.GetRunLogs returned error: %v", err)
	}

	if logs != wantLogs {
		t.Errorf("Expected logs %s, got %s", wantLogs, logs)
	}
}

func TestPipelinesService_RetryRun(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	want := &PipelineRun{
		ID:         "run-2",
		PipelineID: "pipeline-1",
		State:      "pending",
	}

	mux.HandleFunc("/pipelines/pipeline-1/runs/run-1/retry", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	ctx := context.Background()
	run, _, err := client.Pipelines.RetryRun(ctx, "pipeline-1", "run-1")

	if err != nil {
		t.Fatalf("Pipelines.RetryRun returned error: %v", err)
	}

	if run.ID != want.ID {
		t.Errorf("Expected run ID %s, got %s", want.ID, run.ID)
	}
}

// Additional error handling tests

func TestErrorResponse_OtherErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		checkFunc  func(error) bool
	}{
		{"BadRequest", http.StatusBadRequest, IsBadRequest},
		{"Unauthorized", http.StatusUnauthorized, IsUnauthorized},
		{"Forbidden", http.StatusForbidden, IsForbidden},
		{"Conflict", http.StatusConflict, IsConflict},
		{"RateLimited", http.StatusTooManyRequests, IsRateLimited},
		{"ServerError", http.StatusInternalServerError, IsServerError},
		{"UnprocessableEntity", http.StatusUnprocessableEntity, IsUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &ErrorResponse{
				Response: &http.Response{
					StatusCode: tt.statusCode,
				},
			}

			if !tt.checkFunc(err) {
				t.Errorf("%s check failed for status code %d", tt.name, tt.statusCode)
			}
		})
	}
}

// Additional edge case and error handling tests

func TestClient_NewRequestErrors(t *testing.T) {
	// Test with invalid URL
	client := NewClient(WithAPIKey("test-key"))

	// This should work fine with a normal path
	ctx := context.Background()
	req, err := client.NewRequest(ctx, "GET", "test/path", nil)
	if err != nil {
		t.Fatalf("Expected no error for valid path, got: %v", err)
	}

	// Verify request was created properly
	if req.Method != "GET" {
		t.Errorf("Expected GET method, got %s", req.Method)
	}

	if req.Header.Get("Authorization") != "Bearer test-key" {
		t.Errorf("Expected Authorization header to be set")
	}

	if req.Header.Get("User-Agent") == "" {
		t.Error("Expected User-Agent header to be set")
	}
}

func TestClient_NewRequestWithBody(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))
	ctx := context.Background()

	body := map[string]string{"name": "test"}
	req, err := client.NewRequest(ctx, "POST", "test", body)

	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type header to be set to application/json")
	}
}

func TestClient_NewRequestWithQueryParams(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))
	ctx := context.Background()

	opts := &ListOptions{
		Page:    2,
		PerPage: 50,
	}

	req, err := client.NewRequestWithQuery(ctx, "GET", "test", opts)

	if err != nil {
		t.Fatalf("NewRequestWithQuery returned error: %v", err)
	}

	// Check query parameters
	query := req.URL.Query()
	if query.Get("page") != "2" {
		t.Errorf("Expected page=2, got page=%s", query.Get("page"))
	}

	if query.Get("per_page") != "50" {
		t.Errorf("Expected per_page=50, got per_page=%s", query.Get("per_page"))
	}
}

func TestClient_DoWithNoContent(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, "DELETE", "test", nil)

	resp, err := client.Do(req, nil)

	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestClient_DoWithIOWriter(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	testContent := "test file content"

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(testContent)); err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	})

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, "GET", "test", nil)

	var buf bytes.Buffer
	_, err := client.Do(req, &buf)

	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}

	if buf.String() != testContent {
		t.Errorf("Expected content %s, got %s", testContent, buf.String())
	}
}

func TestCheckResponse_Success(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}

	err := CheckResponse(resp)
	if err != nil {
		t.Errorf("CheckResponse returned error for 200 status: %v", err)
	}
}

func TestCheckResponse_WithErrorBody(t *testing.T) {
	errorJSON := `{"message": "Not found", "code": "NOT_FOUND"}`
	resp := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewBufferString(errorJSON)),
	}

	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("Expected ErrorResponse, got %T", err)
	}

	if errResp.Message != "Not found" {
		t.Errorf("Expected message 'Not found', got %s", errResp.Message)
	}

	if errResp.Code != "NOT_FOUND" {
		t.Errorf("Expected code 'NOT_FOUND', got %s", errResp.Code)
	}
}

func TestCheckResponse_WithPlainTextError(t *testing.T) {
	errorText := "Internal Server Error"
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewBufferString(errorText)),
	}

	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("Expected ErrorResponse, got %T", err)
	}

	if errResp.Message != errorText {
		t.Errorf("Expected message '%s', got %s", errorText, errResp.Message)
	}
}

func TestIsClientError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"400 is client error", http.StatusBadRequest, true},
		{"404 is client error", http.StatusNotFound, true},
		{"499 is client error", 499, true},
		{"500 is not client error", http.StatusInternalServerError, false},
		{"200 is not client error", http.StatusOK, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &ErrorResponse{
				Response: &http.Response{
					StatusCode: tt.statusCode,
				},
			}

			result := IsClientError(err)
			if result != tt.expected {
				t.Errorf("IsClientError(%d) = %v, want %v", tt.statusCode, result, tt.expected)
			}
		})
	}
}

func TestPopulatePageValues(t *testing.T) {
	tests := []struct {
		name       string
		linkHeader string
		wantNext   int
		wantPrev   int
		wantFirst  int
		wantLast   int
	}{
		{
			name:       "full pagination",
			linkHeader: `<http://api.test/items?page=3>; rel="next", <http://api.test/items?page=1>; rel="prev", <http://api.test/items?page=1>; rel="first", <http://api.test/items?page=10>; rel="last"`,
			wantNext:   3,
			wantPrev:   1,
			wantFirst:  1,
			wantLast:   10,
		},
		{
			name:       "only next",
			linkHeader: `<http://api.test/items?page=2>; rel="next"`,
			wantNext:   2,
		},
		{
			name:       "empty header",
			linkHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &Response{
				Response: &http.Response{
					Header: http.Header{},
				},
			}

			if tt.linkHeader != "" {
				resp.Header.Set("Link", tt.linkHeader)
			}

			resp.populatePageValues()

			if resp.NextPage != tt.wantNext {
				t.Errorf("NextPage = %d, want %d", resp.NextPage, tt.wantNext)
			}
			if resp.PrevPage != tt.wantPrev {
				t.Errorf("PrevPage = %d, want %d", resp.PrevPage, tt.wantPrev)
			}
			if resp.FirstPage != tt.wantFirst {
				t.Errorf("FirstPage = %d, want %d", resp.FirstPage, tt.wantFirst)
			}
			if resp.LastPage != tt.wantLast {
				t.Errorf("LastPage = %d, want %d", resp.LastPage, tt.wantLast)
			}
		})
	}
}

func TestErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrorResponse
		contains string
	}{
		{
			name: "with code",
			err: &ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "Resource not found",
				Response: &http.Response{
					StatusCode: 404,
				},
			},
			contains: "NOT_FOUND",
		},
		{
			name: "with request ID",
			err: &ErrorResponse{
				Message:   "Server error",
				RequestID: "req-123",
				Response: &http.Response{
					StatusCode: 500,
				},
			},
			contains: "req-123",
		},
		{
			name: "simple message",
			err: &ErrorResponse{
				Message: "Bad request",
				Response: &http.Response{
					StatusCode: 400,
				},
			},
			contains: "Bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()
			if errStr == "" {
				t.Error("Error() returned empty string")
			}
			if tt.contains != "" && !contains(errStr, tt.contains) {
				t.Errorf("Error message %q does not contain %q", errStr, tt.contains)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Message: "invalid format",
	}

	errStr := err.Error()
	if !contains(errStr, "email") {
		t.Errorf("Expected error to contain field name, got: %s", errStr)
	}
	if !contains(errStr, "invalid format") {
		t.Errorf("Expected error to contain message, got: %s", errStr)
	}
}

func TestRateLimitError(t *testing.T) {
	err := &RateLimitError{
		ErrorResponse: &ErrorResponse{
			Message: "Too many requests",
		},
		RetryAfter: 60,
	}

	errStr := err.Error()
	if !contains(errStr, "60") {
		t.Errorf("Expected error to contain retry after time, got: %s", errStr)
	}
	if !contains(errStr, "rate limited") {
		t.Errorf("Expected error to contain 'rate limited', got: %s", errStr)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
