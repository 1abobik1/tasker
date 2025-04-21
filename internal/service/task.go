package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/1abobik1/tasker/internal/errs"
	"github.com/1abobik1/tasker/internal/models"
	"github.com/google/uuid"
)

type RepositoryI interface {
	Create(ctx context.Context, t models.Task) error
	GetByID(ctx context.Context, id string) (models.Task, error)
	UpdateStatus(ctx context.Context, id string, status models.TaskStatus) error
	SaveResult(ctx context.Context, id string, result []byte) error
	SaveError(ctx context.Context, id, errMsg string) error
}

// публикация в очередь rabbitmq
type MessageBroker interface {
	PublishTask(task models.Task) error
}

type Service struct {
	repo     RepositoryI
	producer MessageBroker
}

func NewService(repo RepositoryI, producer MessageBroker) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
	}
}

func (s *Service) CreateTask(ctx context.Context, taskType string, payload []byte) (string, time.Time, error) {
	op := "location internal.service.CreateTask"

	uuid := uuid.NewString()
	now := time.Now().UTC()

	task := models.Task{
		ID:        uuid,
		Type:      taskType,
		Payload:   payload,
		Status:    models.StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		log.Printf("%v:%v", op, err)
		return "", time.Time{}, errs.ErrInternalServer
	}

	if err := s.producer.PublishTask(task); err != nil {
		log.Printf("%v:%v", op, err)
		return "", time.Time{}, errs.ErrInternalServer
	}

	return uuid, now, nil
}

func (s *Service) GetTask(ctx context.Context, uuid string) (models.Task, error) {
	op := "location internal.service.GetTask"

	task, err := s.repo.GetByID(ctx, uuid)
	if err != nil {
		log.Printf("%v:%v", op, err)
		if errors.Is(err, errs.ErrIDNotFound) {
			return models.Task{}, errs.ErrIDNotFound
		} else {
			return models.Task{}, errs.ErrInternalServer
		}
	}
	return task, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status models.TaskStatus) error {
	op := "location internal.service.UpdateStatus"

	err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Printf("%v:%v", op, err)
		return errs.ErrInternalServer
	}
	return nil
}

func (s *Service) SaveResult(ctx context.Context, id string, result []byte) error {
	op := "location internal.service.SaveResult"

	err := s.repo.SaveResult(ctx, id, result)
	if err != nil {
		log.Printf("%v:%v", op, err)
		return errs.ErrInternalServer
	}
	return nil
}

func (s *Service) SaveError(ctx context.Context, id, errMsg string) error {
	op := "location internal.service.SaveError"

	err := s.repo.SaveError(ctx, id, errMsg)
	if err != nil {
		log.Printf("%v:%v", op, err)
		return errs.ErrInternalServer
	}
	return nil
}
