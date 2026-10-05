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
	"game-event-service/internal/repository"
	"game-event-service/internal/service"
	"game-event-service/internal/worker"
)

func main() {
	cfg := config.Load()

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

	playerCache := cache.NewPlayerCache(redisClient.Client())

	eventRepository := repository.NewEventRepository(db)
	playerRepository := repository.NewPlayerRepository(db)

	eventWorker := worker.NewEventWorker(100, 3)

	go eventWorker.Run(ctx)

	eventService := service.NewEventService(eventRepository, eventWorker)
	playerService := service.NewPlayerService(playerRepository, playerCache)

	h := handler.New(eventService, playerService)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("server started on :%s", cfg.Port)
		serverErrors <- server.ListenAndServe()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	case sig := <-shutdownSignal:
		log.Printf("shutdown signal received: %s", sig)
	}

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}

	log.Println("server stopped")
}
