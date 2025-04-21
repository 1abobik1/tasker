package dto

import (
	"time"

	"github.com/1abobik1/tasker/internal/models"
)

type CreateTaskRequest struct {
	Type    string `json:"type"  binding:"required"`   // пример, "resize_image", "fetch_url"
	Payload string `json:"payload" binding:"required"` // имитация задачи
}

type CreateTaskResponse struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Status    models.TaskStatus `json:"status"`
	CreatedAt time.Time         `json:"createdAt"`
}

type GetTaskResponse struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Status    models.TaskStatus `json:"status"`
	Result    interface{}       `json:"result,omitempty"`
	Error     string            `json:"error,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}
