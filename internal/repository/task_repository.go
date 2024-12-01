package repository

import (
    "context"
    "knovel/internal/models"
)

type TaskRepository interface {
    Create(ctx context.Context, task *models.Task) error
    Update(ctx context.Context, task *models.Task) error
    GetByID(ctx context.Context, id uint) (*models.Task, error)
    GetTasksByAssignee(ctx context.Context, assigneeID uint) ([]models.Task, error)
    GetFilteredTasks(ctx context.Context, filter models.TaskFilter) ([]models.Task, error)
    GetEmployeeTaskSummary(ctx context.Context) ([]EmployeeTaskSummary, error)
}

type EmployeeTaskSummary struct {
    EmployeeID   uint   `json:"employee_id"`
    EmployeeName string `json:"employee_name"`
    TotalTasks   int    `json:"total_tasks"`
    Completed    int    `json:"completed_tasks"`
} 