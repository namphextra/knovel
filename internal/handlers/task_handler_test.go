package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"knovel/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockTaskRepository)
	handler := NewTaskHandler(mockRepo)

	// Test case: successful task creation
	t.Run("successful task creation", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Task")).Return(nil)

		task := models.Task{
			Title:       "Test Task",
			Description: "Test Description",
			AssigneeID:  1,
		}

		body, _ := json.Marshal(task)
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Mock authenticated user
		c.Set("user", &models.User{ID: 1, Role: models.RoleEmployer})

		handler.CreateTask(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
