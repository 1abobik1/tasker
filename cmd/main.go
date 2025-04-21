package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1abobik1/tasker/config"
	"github.com/1abobik1/tasker/internal/broker/rabbitmq"
	"github.com/1abobik1/tasker/internal/db"
	"github.com/1abobik1/tasker/internal/handler"
	"github.com/1abobik1/tasker/internal/repository"
	"github.com/1abobik1/tasker/internal/service"
	"github.com/1abobik1/tasker/internal/worker"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.MustLoad()

	// init БД
	pgx := db.InitPostgres(cfg.PostgresURL)
	defer pgx.Close()

	//  подключение к RabbitMQ
	var conn *amqp.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(cfg.RabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ dial failed (attempt %d): %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	queueName := cfg.QueueName
	if _, err := ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	); err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	producer := rabbitmq.NewProducer(ch, queueName)
	repo := repository.NewPostgresRepo(pgx)
	svc := service.NewService(repo, producer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	w := worker.NewWorker(ch, queueName, svc)
	go func() {
		if err := w.Start(ctx); err != nil {
			log.Printf("Worker error: %v", err)
			cancel()
		}
	}()

	router := gin.Default()
	h := handler.NewHandler(svc)
	router.POST("/task", h.CreateTask)
	router.GET("/task", h.GetTask)

	srvErr := make(chan error, 1)
	go func() {
		log.Printf("HTTP server listening on %s", cfg.HTTPPort)
		srvErr <- router.Run(cfg.HTTPPort)
	}()

	select {
	case sig := <-sigCh:
		log.Printf("Received signal %s, shutting down...", sig)
	case err := <-srvErr:
		log.Printf("HTTP server error: %v", err)
	}

	cancel()

	time.Sleep(4 * time.Second)
	log.Println("Server exited gracefully")
}
