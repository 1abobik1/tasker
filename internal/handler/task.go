package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/1abobik1/tasker/internal/dto"
	"github.com/1abobik1/tasker/internal/errs"
	"github.com/1abobik1/tasker/internal/models"
	"github.com/gin-gonic/gin"
)

type ServiceI interface {
	CreateTask(ctx context.Context, input []byte) (string, time.Time, error)
	GetTask(ctx context.Context, uuid string) (models.Task, error)
}

type Handler struct {
	svc ServiceI
}

func NewHandler(svc ServiceI) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid, createdAt, err := h.svc.CreateTask(c, []byte(req.Input))
	if err != nil {
		if errors.Is(err, errs.ErrInternalServer) {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusAccepted, dto.CreateTaskResponse{
		ID:        uuid,
		Status:    models.StatusPending,
		CreatedAt: createdAt,
	})
}

func (h *Handler) GetTask(c *gin.Context) {
	uuid := c.Query("id")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be specified in query parameters."})
		return
	}
	task, err := h.svc.GetTask(c, uuid)
	if err != nil {
		if errors.Is(err, errs.ErrIDNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found" })
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	resp := dto.GetTaskResponse{
		ID:        task.ID,
		Status:    task.Status,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	if task.Status == models.StatusCompleted {
		resp.Result = task.Result
	}
	if task.Error.Valid {
		resp.Error = task.Error.String
		resp.Result = nil
	}

	c.JSON(http.StatusOK, resp)
}
