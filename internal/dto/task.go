package dto

import (
	"time"

	"github.com/1abobik1/tasker/internal/models"
)

type CreateTaskRequest struct {
	Input string `json:"input" binding:"required"`
}

type CreateTaskResponse struct {
	ID        string            `json:"id"`
	Status    models.TaskStatus `json:"status"`
	CreatedAt time.Time         `json:"createdAt"`
}

type GetTaskResponse struct {
	ID        string            `json:"id"`
	Status    models.TaskStatus `json:"status"`
	Result    interface{}       `json:"result,omitempty"`
	Error     string            `json:"error,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}
