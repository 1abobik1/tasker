package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/1abobik1/tasker/internal/models"
	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	ch    *amqp091.Channel
	queue string
}

func NewProducer(ch *amqp091.Channel, queue string) *Producer {
	return &Producer{ch: ch, queue: queue}
}

func (p *Producer) PublishTask(task models.Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	err = p.ch.Publish(
		"",         
		p.queue,    
		false,      
		false,      
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish task: %w", err)
	}

	return nil
}
