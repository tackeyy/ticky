package ticktick

import "time"

// User represents a TickTick user.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// Project represents a TickTick project (list).
type Project struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color,omitempty"`
	SortOrder  int64  `json:"sortOrder"`
	ViewMode   string `json:"viewMode,omitempty"`
	Kind       string `json:"kind,omitempty"`
	GroupID    string `json:"groupId,omitempty"`
	IsOwner    bool   `json:"isOwner"`
	Permission string `json:"permission,omitempty"`
}

// Task represents a TickTick task.
type Task struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"projectId"`
	Title       string    `json:"title"`
	Content     string    `json:"content,omitempty"`
	Desc        string    `json:"desc,omitempty"`
	Priority    int       `json:"priority"`
	Status      int       `json:"status"`
	DueDate     string    `json:"dueDate,omitempty"`
	StartDate   string    `json:"startDate,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	TimeZone    string    `json:"timeZone,omitempty"`
	IsAllDay    bool      `json:"isAllDay"`
	CompletedAt string    `json:"completedTime,omitempty"`
	CreatedAt   time.Time `json:"createdTime,omitempty"`
	ModifiedAt  time.Time `json:"modifiedTime,omitempty"`
}

// ProjectData is the response from GET /project/{id}/data.
type ProjectData struct {
	Project Project `json:"project"`
	Tasks   []Task  `json:"tasks"`
}

// TaskCreateRequest is the request body for creating a task.
type TaskCreateRequest struct {
	Title     string   `json:"title"`
	ProjectID string   `json:"projectId,omitempty"`
	Content   string   `json:"content,omitempty"`
	Desc      string   `json:"desc,omitempty"`
	Priority  int      `json:"priority,omitempty"`
	DueDate   string   `json:"dueDate,omitempty"`
	StartDate string   `json:"startDate,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	TimeZone  string   `json:"timeZone,omitempty"`
	IsAllDay  bool     `json:"isAllDay,omitempty"`
}

// TaskUpdateRequest is the request body for updating a task.
type TaskUpdateRequest struct {
	ID        string   `json:"id"`
	ProjectID string   `json:"projectId"`
	Title     string   `json:"title,omitempty"`
	Content   string   `json:"content,omitempty"`
	Desc      string   `json:"desc,omitempty"`
	Priority  *int     `json:"priority,omitempty"`
	DueDate   *string  `json:"dueDate"`
	StartDate *string  `json:"startDate,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	TimeZone  string   `json:"timeZone,omitempty"`
	IsAllDay  *bool    `json:"isAllDay,omitempty"`
}

// OAuthToken represents the OAuth token response.
type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
}
