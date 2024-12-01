package repository

import (
    "context"
    "knovel/internal/models"
    "gorm.io/gorm"
)

type GormTaskRepository struct {
    db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) *GormTaskRepository {
    return &GormTaskRepository{db: db}
}

func (r *GormTaskRepository) Create(ctx context.Context, task *models.Task) error {
    return r.db.WithContext(ctx).Create(task).Error
}

func (r *GormTaskRepository) Update(ctx context.Context, task *models.Task) error {
    return r.db.WithContext(ctx).Save(task).Error
}

func (r *GormTaskRepository) GetByID(ctx context.Context, id uint) (*models.Task, error) {
    var task models.Task
    err := r.db.WithContext(ctx).First(&task, id).Error
    return &task, err
}

func (r *GormTaskRepository) GetTasksByAssignee(ctx context.Context, assigneeID uint) ([]models.Task, error) {
    var tasks []models.Task
    err := r.db.WithContext(ctx).Where("assignee_id = ?", assigneeID).Find(&tasks).Error
    return tasks, err
}

func (r *GormTaskRepository) GetFilteredTasks(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
    query := r.db.WithContext(ctx)

    if filter.AssigneeID != nil {
        query = query.Where("assignee_id = ?", *filter.AssigneeID)
    }
    
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }

    if filter.SortBy != "" {
        order := "asc"
        if filter.SortOrder == "desc" {
            order = "desc"
        }
        query = query.Order(filter.SortBy + " " + order)
    }

    var tasks []models.Task
    err := query.Find(&tasks).Error
    return tasks, err
}

func (r *GormTaskRepository) GetEmployeeTaskSummary(ctx context.Context) ([]EmployeeTaskSummary, error) {
    var summaries []EmployeeTaskSummary
    
    err := r.db.WithContext(ctx).
        Table("users").
        Select(`
            users.id as employee_id,
            users.name as employee_name,
            COUNT(tasks.id) as total_tasks,
            COUNT(CASE WHEN tasks.status = ? THEN 1 END) as completed_tasks
        `, models.TaskStatusCompleted).
        Joins("LEFT JOIN tasks ON tasks.assignee_id = users.id").
        Where("users.role = ?", models.RoleEmployee).
        Group("users.id, users.name").
        Find(&summaries).Error

    return summaries, err
} 