package models

import (
	"database/sql"
	"time"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID        string
	Payload   []byte
	Status    TaskStatus
	Result    []byte
	Error     sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}
