package worker

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	"game-event-service/internal/metrics"
	"game-event-service/internal/model"
)

func BenchmarkEventWorkerSubmit(b *testing.B) {
	originalOutput := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(originalOutput)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appMetrics := metrics.New()
	eventWorker := NewEventWorker(100, 4, appMetrics)

	go eventWorker.Run(ctx)

	event := model.Event{
		PlayerID: 1,
		Type:     "purchase.created",
		Payload:  "benchmark",
	}

	b.ResetTimer()

	for b.Loop() {
		eventWorker.Submit(event)
	}

	b.StopTimer()

	cancel()
	time.Sleep(10 * time.Millisecond)
}
