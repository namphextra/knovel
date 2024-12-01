package handlers

import (
	"knovel/internal/models"
	"knovel/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskRepo repository.TaskRepository
}

func NewTaskHandler(taskRepo repository.TaskRepository) *TaskHandler {
	return &TaskHandler{taskRepo: taskRepo}
}

// CreateTask - Only employers can create tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	if user.Role != models.RoleEmployer {
		c.JSON(http.StatusForbidden, gin.H{"error": "only employers can create tasks"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.CreatedBy = user.ID
	if err := h.taskRepo.Create(c, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTasks - Returns tasks based on user role
func (h *TaskHandler) GetTasks(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var tasks []models.Task
	var err error

	if user.Role == models.RoleEmployee {
		// Employees can only see their tasks
		tasks, err = h.taskRepo.GetTasksByAssignee(c, user.ID)
	} else {
		// Employers can see all tasks with filtering
		var filter models.TaskFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tasks, err = h.taskRepo.GetFilteredTasks(c, filter)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetEmployeeTaskSummary(c *gin.Context) {
	// Implement the logic to get the employee task summary
	c.JSON(http.StatusOK, gin.H{"message": "Task summary for employer"})
}
