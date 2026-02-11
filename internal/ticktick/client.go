package ticktick

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const baseURL = "https://api.ticktick.com/open/v1"

// Client is the TickTick API client.
type Client struct {
	httpClient  *http.Client
	accessToken string
}

// NewClient creates a new TickTick client.
// It checks TICKTICK_ACCESS_TOKEN env var first, then falls back to token file.
func NewClient() (*Client, error) {
	if token := os.Getenv("TICKTICK_ACCESS_TOKEN"); token != "" {
		return &Client{
			httpClient:  &http.Client{Timeout: 30 * time.Second},
			accessToken: token,
		}, nil
	}

	savedToken, err := LoadToken()
	if err != nil {
		return nil, fmt.Errorf("not authenticated: run 'ticky auth login' first (%w)", err)
	}

	// Refresh if expired
	if savedToken.ExpiresAt > 0 && time.Now().Unix() > savedToken.ExpiresAt {
		if savedToken.RefreshToken == "" {
			return nil, fmt.Errorf("token expired and no refresh token available: run 'ticky auth login'")
		}
		newToken, err := RefreshAccessToken(savedToken.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w (run 'ticky auth login')", err)
		}
		if err := SaveToken(newToken); err != nil {
			return nil, fmt.Errorf("failed to save refreshed token: %w", err)
		}
		savedToken = newToken
	}

	return &Client{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		accessToken: savedToken.AccessToken,
	}, nil
}

// Get performs a GET request.
func (c *Client) Get(path string) ([]byte, error) {
	return c.do("GET", path, nil)
}

// Post performs a POST request with JSON body.
func (c *Client) Post(path string, body any) ([]byte, error) {
	if body == nil {
		return c.do("POST", path, nil)
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	return c.do("POST", path, bytes.NewBuffer(data))
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) ([]byte, error) {
	return c.do("DELETE", path, nil)
}

func (c *Client) do(method, path string, body io.Reader) ([]byte, error) {
	url := baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetProjects returns all projects.
func (c *Client) GetProjects() ([]Project, error) {
	data, err := c.Get("/project")
	if err != nil {
		return nil, err
	}
	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}
	return projects, nil
}

// GetProject returns a single project by ID.
func (c *Client) GetProject(id string) (*Project, error) {
	data, err := c.Get("/project/" + id)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}
	return &project, nil
}

// GetProjectData returns a project with its tasks.
func (c *Client) GetProjectData(projectID string) (*ProjectData, error) {
	data, err := c.Get("/project/" + projectID + "/data")
	if err != nil {
		return nil, err
	}
	var pd ProjectData
	if err := json.Unmarshal(data, &pd); err != nil {
		return nil, fmt.Errorf("failed to parse project data: %w", err)
	}
	return &pd, nil
}

// CreateTask creates a new task.
func (c *Client) CreateTask(req *TaskCreateRequest) (*Task, error) {
	data, err := c.Post("/task", req)
	if err != nil {
		return nil, err
	}
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}
	return &task, nil
}

// GetTask returns a single task by project ID and task ID.
func (c *Client) GetTask(projectID, taskID string) (*Task, error) {
	data, err := c.Get("/project/" + projectID + "/task/" + taskID)
	if err != nil {
		return nil, err
	}
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}
	return &task, nil
}

// UpdateTask updates an existing task.
func (c *Client) UpdateTask(req *TaskUpdateRequest) (*Task, error) {
	data, err := c.Post("/task/"+req.ID, req)
	if err != nil {
		return nil, err
	}
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}
	return &task, nil
}

// CompleteTask marks a task as complete.
func (c *Client) CompleteTask(projectID, taskID string) error {
	_, err := c.Post("/project/"+projectID+"/task/"+taskID+"/complete", nil)
	return err
}

// DeleteTask deletes a task.
func (c *Client) DeleteTask(projectID, taskID string) error {
	_, err := c.Delete("/project/" + projectID + "/task/" + taskID)
	return err
}

// DiscoverInboxID discovers the Inbox project ID by creating and deleting a temporary task.
func (c *Client) DiscoverInboxID() (string, error) {
	// Check cache first
	if id, err := LoadInboxID(); err == nil && id != "" {
		return id, nil
	}

	task, err := c.CreateTask(&TaskCreateRequest{Title: ".ticky-inbox-probe"})
	if err != nil {
		return "", fmt.Errorf("failed to discover inbox ID: %w", err)
	}
	inboxID := task.ProjectID
	_ = c.DeleteTask(inboxID, task.ID)

	// Cache for future use
	_ = SaveInboxID(inboxID)

	return inboxID, nil
}

// GetAllProjectIDs returns all project IDs including Inbox.
func (c *Client) GetAllProjectIDs() ([]string, error) {
	projects, err := c.GetProjects()
	if err != nil {
		return nil, err
	}

	inboxID, err := c.DiscoverInboxID()
	if err != nil {
		return nil, err
	}

	ids := []string{inboxID}
	for _, p := range projects {
		if p.ID != inboxID {
			ids = append(ids, p.ID)
		}
	}

	return ids, nil
}
