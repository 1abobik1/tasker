// internal/worker/worker.go

package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/1abobik1/tasker/internal/models"
	"github.com/1abobik1/tasker/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	ch    *amqp.Channel
	queue string
	svc   *service.Service
	registry *Registry
}

func NewWorker(ch *amqp.Channel, queue string, svc *service.Service, registry *Registry) *Worker {
	return &Worker{ch: ch, queue: queue, svc: svc, registry: registry}
}

func (w *Worker) Start(ctx context.Context) error {
	msgs, err := w.ch.Consume(
		w.queue, "", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	log.Println("Worker started. Waiting for messages...")

	for msg := range msgs {
		go w.handleMessage(ctx, msg.Body)
	}

	return nil
}

func (w *Worker) handleMessage(ctx context.Context, body []byte) {
	var task models.Task
	if err := json.Unmarshal(body, &task); err != nil {
		log.Printf("failed to unmarshal task: %v", err)
		return
	}

	processor := w.registry.GetProcessor(task.Type)
	if processor == nil {
		log.Printf("unknown task type: %s", task.Type)
		_ = w.svc.SaveError(ctx, task.ID, "unknown task type: "+task.Type)
		_ = w.svc.UpdateStatus(ctx, task.ID, models.StatusFailed)
		return
	}

	result, err := processor.Process(ctx, task.Payload)
	if err != nil {
		log.Printf("failed to process task %s: %v", task.ID, err)
		_ = w.svc.SaveError(ctx, task.ID, err.Error())
		_ = w.svc.UpdateStatus(ctx, task.ID, models.StatusFailed)
		return
	}

	if err := w.svc.SaveResult(ctx, task.ID, result); err != nil {
		log.Printf("failed to save result: %v", err)
		return
	}

	_ = w.svc.UpdateStatus(ctx, task.ID, models.StatusCompleted)
}
