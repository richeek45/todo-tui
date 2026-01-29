package models

import "time"

type Status string
type Priority string

const (
	Completed  Status = "completed"
	NotStarted Status = "pending"
	InProgress Status = "in_progress"
)

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	Error        error
	StartTime    time.Time
	FinishedTime *time.Time

	Id          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TodoFilter struct {
	Status        Status     `json:"status,omitempty"`
	Priority      Priority   `json:"priority,omitempty"`
	FromDate      *time.Time `json:"from_date,omitempty"`
	ToDate        *time.Time `json:"to_date,omitempty"`
	CompletedFrom *time.Time `json:"completed_from,omitempty"`
	CompletedTo   *time.Time `json:"completed_to,omitempty"`
	Search        string     `json:"search,omitempty"`
}

type CursorPagination struct {
	Cursor   string `json:"cursor,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	OrderBy  string `json:"order_by,omitempty"`
	OrderDir string `json:"order_dir,omitempty"`
}

type PaginatedTodos struct {
	Todos      []Task `json:"todos"`
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
	TotalCount int    `json:"total_count,omitempty"`
}
