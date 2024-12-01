package models

import (
	"fmt"
	"time"
)

func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("title is required")
	}
	if t.AssigneeID == 0 {
		return fmt.Errorf("assignee is required")
	}
	if t.DueDate.Before(time.Now()) {
		return fmt.Errorf("due date must be in the future")
	}
	return nil
} 