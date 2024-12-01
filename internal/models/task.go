package models

import "time"

type TaskStatus string

const (
    TaskStatusPending    TaskStatus = "pending"
    TaskStatusInProgress TaskStatus = "in_progress"
    TaskStatusCompleted  TaskStatus = "completed"
)

type Task struct {
    ID          uint       `json:"id" gorm:"primaryKey"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      TaskStatus `json:"status"`
    AssigneeID  uint       `json:"assignee_id"`
    CreatedBy   uint       `json:"created_by"`
    DueDate     time.Time  `json:"due_date"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskFilter struct {
    AssigneeID *uint      `json:"assignee_id,omitempty"`
    Status     TaskStatus `json:"status,omitempty"`
    SortBy     string     `json:"sort_by,omitempty"` // "due_date", "status", "created_at"
    SortOrder  string     `json:"sort_order,omitempty"` // "asc", "desc"
} 