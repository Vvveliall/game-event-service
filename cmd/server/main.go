package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"game-event-service/internal/cache"
	"game-event-service/internal/config"
	"game-event-service/internal/database"
	"game-event-service/internal/handler"
	"game-event-service/internal/logger"
	"game-event-service/internal/metrics"
	"game-event-service/internal/middleware"
	"game-event-service/internal/queue"
	"game-event-service/internal/repository"
	"game-event-service/internal/service"
	"game-event-service/internal/worker"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New(cfg.Env)
	appMetrics := metrics.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	redisClient, err := cache.NewRedisClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	defer redisClient.Close()

	rabbitMQ, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("rabbitmq connection failed: %v", err)
	}
	defer rabbitMQ.Close()

	playerCache := cache.NewPlayerCache(redisClient.Client())

	eventRepository := repository.NewEventRepository(db)
	playerRepository := repository.NewPlayerRepository(db)

	eventWorker := worker.NewEventWorker(100, 3, appMetrics)
	eventConsumer := queue.NewEventConsumer(rabbitMQ, eventWorker)

	go eventWorker.Run(ctx)

	go func() {
		if err := eventConsumer.Run(ctx); err != nil {
			appLogger.Error("event consumer error", "error", err)
			cancel()
		}
	}()

	eventService := service.NewEventService(eventRepository, rabbitMQ)
	playerService := service.NewPlayerService(playerRepository, playerCache)

	h := handler.New(eventService, playerService)

	routes := h.Routes()
	routes = middleware.Metrics(appMetrics)(routes)
	routes = middleware.Logging(appLogger)(routes)
	routes = middleware.RequestID(routes)

	mux := http.NewServeMux()
	mux.Handle("/", routes)
	mux.Handle("/metrics", appMetrics.Handler())

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		appLogger.Info("server started", "port", cfg.Port)

		if err := server.ListenAndServe(); err != nil {
			serverErrors <- err
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			appLogger.Error("server error", "error", err)
		}
	case sig := <-shutdownSignal:
		appLogger.Info("shutdown signal received", "signal", sig.String())
	}

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("graceful shutdown error", "error", err)
	}

	appLogger.Info("server stopped")
}
