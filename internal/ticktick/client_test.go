package ticktick_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tackeyy/ticky/internal/ticktick"
)

// setupMockServer creates a mock HTTP server and a Client pointing to it.
// The returned cleanup function restores the original baseURL and closes the server.
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*ticktick.Client, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	restore := ticktick.SetBaseURL(server.URL)
	client := ticktick.NewTestClient(server.Client(), "test-token")
	return client, func() {
		restore()
		server.Close()
	}
}

// --- Project Operations ---

func TestGetProjects_Success(t *testing.T) {
	// Arrange
	want := []ticktick.Project{
		{ID: "proj-1", Name: "Work"},
		{ID: "proj-2", Name: "Personal"},
	}
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/project" {
			t.Errorf("path = %s, want /project", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	})
	defer cleanup()

	// Act
	got, err := client.GetProjects()

	// Assert
	if err != nil {
		t.Fatalf("GetProjects() returned unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("GetProjects() returned %d projects, want 2", len(got))
	}
	if got[0].ID != "proj-1" || got[0].Name != "Work" {
		t.Errorf("GetProjects()[0] = %+v, want ID=proj-1, Name=Work", got[0])
	}
	if got[1].ID != "proj-2" || got[1].Name != "Personal" {
		t.Errorf("GetProjects()[1] = %+v, want ID=proj-2, Name=Personal", got[1])
	}
}

func TestGetProjects_HTTPError(t *testing.T) {
	// Arrange
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	})
	defer cleanup()

	// Act
	_, err := client.GetProjects()

	// Assert
	if err == nil {
		t.Fatal("GetProjects() expected error for 500 response, got nil")
	}
	if !strings.Contains(err.Error(), "status 500") {
		t.Errorf("error = %q, want to contain 'status 500'", err.Error())
	}
}

func TestGetProject_Success(t *testing.T) {
	// Arrange
	want := ticktick.Project{ID: "proj-1", Name: "Work", IsOwner: true}
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/project/proj-1" {
			t.Errorf("path = %s, want /project/proj-1", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	})
	defer cleanup()

	// Act
	got, err := client.GetProject("proj-1")

	// Assert
	if err != nil {
		t.Fatalf("GetProject() returned unexpected error: %v", err)
	}
	if got.ID != "proj-1" {
		t.Errorf("GetProject().ID = %q, want %q", got.ID, "proj-1")
	}
	if got.Name != "Work" {
		t.Errorf("GetProject().Name = %q, want %q", got.Name, "Work")
	}
	if !got.IsOwner {
		t.Error("GetProject().IsOwner = false, want true")
	}
}

func TestGetProjectData_Success(t *testing.T) {
	// Arrange
	want := ticktick.ProjectData{
		Project: ticktick.Project{ID: "proj-1", Name: "Work"},
		Tasks: []ticktick.Task{
			{ID: "task-1", ProjectID: "proj-1", Title: "Task One"},
			{ID: "task-2", ProjectID: "proj-1", Title: "Task Two"},
		},
	}
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/project/proj-1/data" {
			t.Errorf("path = %s, want /project/proj-1/data", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	})
	defer cleanup()

	// Act
	got, err := client.GetProjectData("proj-1")

	// Assert
	if err != nil {
		t.Fatalf("GetProjectData() returned unexpected error: %v", err)
	}
	if got.Project.ID != "proj-1" {
		t.Errorf("GetProjectData().Project.ID = %q, want %q", got.Project.ID, "proj-1")
	}
	if len(got.Tasks) != 2 {
		t.Fatalf("GetProjectData().Tasks has %d items, want 2", len(got.Tasks))
	}
	if got.Tasks[0].Title != "Task One" {
		t.Errorf("GetProjectData().Tasks[0].Title = %q, want %q", got.Tasks[0].Title, "Task One")
	}
}

// --- Task Operations ---

func TestCreateTask_Success(t *testing.T) {
	// Arrange
	req := &ticktick.TaskCreateRequest{
		Title:     "New Task",
		ProjectID: "proj-1",
		Priority:  3,
	}
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/task" {
			t.Errorf("path = %s, want /task", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		// Verify request body
		body, _ := io.ReadAll(r.Body)
		var reqBody ticktick.TaskCreateRequest
		if err := json.Unmarshal(body, &reqBody); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if reqBody.Title != "New Task" {
			t.Errorf("request body Title = %q, want %q", reqBody.Title, "New Task")
		}
		if reqBody.ProjectID != "proj-1" {
			t.Errorf("request body ProjectID = %q, want %q", reqBody.ProjectID, "proj-1")
		}
		if reqBody.Priority != 3 {
			t.Errorf("request body Priority = %d, want %d", reqBody.Priority, 3)
		}

		resp := ticktick.Task{ID: "task-new", ProjectID: "proj-1", Title: "New Task", Priority: 3}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer cleanup()

	// Act
	got, err := client.CreateTask(req)

	// Assert
	if err != nil {
		t.Fatalf("CreateTask() returned unexpected error: %v", err)
	}
	if got.ID != "task-new" {
		t.Errorf("CreateTask().ID = %q, want %q", got.ID, "task-new")
	}
	if got.Title != "New Task" {
		t.Errorf("CreateTask().Title = %q, want %q", got.Title, "New Task")
	}
}

func TestGetTask_Success(t *testing.T) {
	// Arrange
	want := ticktick.Task{ID: "task-1", ProjectID: "proj-1", Title: "My Task"}
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/project/proj-1/task/task-1" {
			t.Errorf("path = %s, want /project/proj-1/task/task-1", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	})
	defer cleanup()

	// Act
	got, err := client.GetTask("proj-1", "task-1")

	// Assert
	if err != nil {
		t.Fatalf("GetTask() returned unexpected error: %v", err)
	}
	if got.ID != "task-1" {
		t.Errorf("GetTask().ID = %q, want %q", got.ID, "task-1")
	}
	if got.Title != "My Task" {
		t.Errorf("GetTask().Title = %q, want %q", got.Title, "My Task")
	}
}

func TestUpdateTask_Success(t *testing.T) {
	// Arrange
	newTitle := "Updated Title"
	req := &ticktick.TaskUpdateRequest{
		ID:        "task-1",
		ProjectID: "proj-1",
		Title:     newTitle,
	}
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/task/task-1" {
			t.Errorf("path = %s, want /task/task-1", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody ticktick.TaskUpdateRequest
		if err := json.Unmarshal(body, &reqBody); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if reqBody.Title != newTitle {
			t.Errorf("request body Title = %q, want %q", reqBody.Title, newTitle)
		}

		resp := ticktick.Task{ID: "task-1", ProjectID: "proj-1", Title: newTitle}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer cleanup()

	// Act
	got, err := client.UpdateTask(req)

	// Assert
	if err != nil {
		t.Fatalf("UpdateTask() returned unexpected error: %v", err)
	}
	if got.Title != newTitle {
		t.Errorf("UpdateTask().Title = %q, want %q", got.Title, newTitle)
	}
}

func TestCompleteTask_Success(t *testing.T) {
	// Arrange
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/project/proj-1/task/task-1/complete" {
			t.Errorf("path = %s, want /project/proj-1/task/task-1/complete", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	// Act
	err := client.CompleteTask("proj-1", "task-1")

	// Assert
	if err != nil {
		t.Fatalf("CompleteTask() returned unexpected error: %v", err)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	// Arrange
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/project/proj-1/task/task-1" {
			t.Errorf("path = %s, want /project/proj-1/task/task-1", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	// Act
	err := client.DeleteTask("proj-1", "task-1")

	// Assert
	if err != nil {
		t.Fatalf("DeleteTask() returned unexpected error: %v", err)
	}
}

// --- Error Handling ---

func TestClient_Unauthorized(t *testing.T) {
	// Arrange
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	})
	defer cleanup()

	// Act
	_, err := client.GetProjects()

	// Assert
	if err == nil {
		t.Fatal("GetProjects() expected error for 401 response, got nil")
	}
	if !strings.Contains(err.Error(), "status 401") {
		t.Errorf("error = %q, want to contain 'status 401'", err.Error())
	}
}

func TestClient_ServerError(t *testing.T) {
	// Arrange
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	})
	defer cleanup()

	// Act
	_, err := client.GetProject("proj-1")

	// Assert
	if err == nil {
		t.Fatal("GetProject() expected error for 500 response, got nil")
	}
	if !strings.Contains(err.Error(), "status 500") {
		t.Errorf("error = %q, want to contain 'status 500'", err.Error())
	}
}

func TestClient_InvalidJSON(t *testing.T) {
	// Arrange
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{invalid json"))
	})
	defer cleanup()

	// Act
	_, err := client.GetProjects()

	// Assert
	if err == nil {
		t.Fatal("GetProjects() expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "parse") {
		t.Errorf("error = %q, want to contain 'parse'", err.Error())
	}
}

// --- Authorization Header ---

func TestClient_AuthorizationHeader(t *testing.T) {
	// Arrange
	var gotAuth string
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	})
	defer cleanup()

	// Act
	_, _ = client.GetProjects()

	// Assert
	want := "Bearer test-token"
	if gotAuth != want {
		t.Errorf("Authorization header = %q, want %q", gotAuth, want)
	}
}

// --- DiscoverInboxID ---

func TestDiscoverInboxID_Success(t *testing.T) {
	// Arrange: isolate HOME so LoadInboxID returns empty (no cache)
	t.Setenv("HOME", t.TempDir())

	callCount := 0
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/task":
			// CreateTask: return a task with projectId = inbox ID
			resp := ticktick.Task{ID: "temp-task-1", ProjectID: "inbox-123", Title: ".ticky-inbox-probe"}
			json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodDelete && r.URL.Path == "/project/inbox-123/task/temp-task-1":
			// DeleteTask: success
			w.WriteHeader(http.StatusOK)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer cleanup()

	// Act
	got, err := client.DiscoverInboxID()

	// Assert
	if err != nil {
		t.Fatalf("DiscoverInboxID() returned unexpected error: %v", err)
	}
	if got != "inbox-123" {
		t.Errorf("DiscoverInboxID() = %q, want %q", got, "inbox-123")
	}
	if callCount != 2 {
		t.Errorf("expected 2 API calls (create + delete), got %d", callCount)
	}
}

func TestDiscoverInboxID_Cached(t *testing.T) {
	// Arrange: set up a cached inbox ID
	t.Setenv("HOME", t.TempDir())

	if err := ticktick.SaveInboxID("cached-inbox-id"); err != nil {
		t.Fatalf("SaveInboxID() returned unexpected error: %v", err)
	}

	apiCalled := false
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		apiCalled = true
		t.Error("API should not be called when inbox ID is cached")
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer cleanup()

	// Act
	got, err := client.DiscoverInboxID()

	// Assert
	if err != nil {
		t.Fatalf("DiscoverInboxID() returned unexpected error: %v", err)
	}
	if got != "cached-inbox-id" {
		t.Errorf("DiscoverInboxID() = %q, want %q", got, "cached-inbox-id")
	}
	if apiCalled {
		t.Error("API was called despite cached inbox ID")
	}
}

// --- GetAllProjectIDs ---

func TestGetAllProjectIDs(t *testing.T) {
	// Arrange: isolate HOME so LoadInboxID returns empty (no cache)
	t.Setenv("HOME", t.TempDir())

	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/project":
			// GetProjects: return projects (one overlaps with inbox)
			projects := []ticktick.Project{
				{ID: "inbox-123", Name: "Inbox"},
				{ID: "proj-1", Name: "Work"},
				{ID: "proj-2", Name: "Personal"},
			}
			json.NewEncoder(w).Encode(projects)
		case r.Method == http.MethodPost && r.URL.Path == "/task":
			// DiscoverInboxID: CreateTask
			resp := ticktick.Task{ID: "temp-task-1", ProjectID: "inbox-123", Title: ".ticky-inbox-probe"}
			json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodDelete && r.URL.Path == "/project/inbox-123/task/temp-task-1":
			// DiscoverInboxID: DeleteTask
			w.WriteHeader(http.StatusOK)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer cleanup()

	// Act
	got, err := client.GetAllProjectIDs()

	// Assert
	if err != nil {
		t.Fatalf("GetAllProjectIDs() returned unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("GetAllProjectIDs() returned %d IDs, want 3", len(got))
	}
	// First element should be the inbox
	if got[0] != "inbox-123" {
		t.Errorf("GetAllProjectIDs()[0] = %q, want %q", got[0], "inbox-123")
	}
	// Check no duplicates
	seen := make(map[string]bool)
	for _, id := range got {
		if seen[id] {
			t.Errorf("GetAllProjectIDs() contains duplicate ID: %q", id)
		}
		seen[id] = true
	}
	// Check all expected IDs are present
	for _, wantID := range []string{"inbox-123", "proj-1", "proj-2"} {
		if !seen[wantID] {
			t.Errorf("GetAllProjectIDs() missing expected ID: %q", wantID)
		}
	}
}
