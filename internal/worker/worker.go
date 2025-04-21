// internal/worker/worker.go

package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/1abobik1/tasker/internal/models"
	"github.com/1abobik1/tasker/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	ch    *amqp.Channel
	queue string
	svc   *service.Service
}

func NewWorker(ch *amqp.Channel, queue string, svc *service.Service) *Worker {
	return &Worker{ch: ch, queue: queue, svc: svc}
}

// Start начинает консюьмить сообщения и обрабатывать их пока не получит ctx.Done
func (w *Worker) Start(ctx context.Context) error {
	msgs, err := w.ch.Consume(
		w.queue, "", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	log.Println("Worker started. Waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker shutting down...")
			return nil
		case d, ok := <-msgs:
			if !ok {
				return nil
			}
			var task models.Task
			if err := json.Unmarshal(d.Body, &task); err != nil {
				log.Printf("Failed to unmarshal task: %v", err)
				continue
			}
			go w.processTask(task)
		}
	}
}

func (w *Worker) processTask(task models.Task) {
	ctx := context.Background()

	log.Printf("Processing task %s...\n", task.ID)

	if err := w.svc.UpdateStatus(ctx, task.ID, models.StatusProcessing); err != nil {
		log.Printf("Failed to set task status to running: %v", err)
		return
	}

	// Simulate work
	time.Sleep(8 * time.Second)

	if string(task.Payload) == "fail" {
		_ = w.svc.SaveError(ctx, task.ID, "simulated error")
		log.Printf("Task %s failed", task.ID)
		return
	}

	result := []byte("result of: " + string(task.Payload))
	if err := w.svc.SaveResult(ctx, task.ID, result); err != nil {
		log.Printf("Failed to save result: %v", err)
		return
	}

	log.Printf("Task %s completed", task.ID)
}
