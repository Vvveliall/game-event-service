package worker

import (
	"context"
	"log"
	"sync"

	"game-event-service/internal/metrics"
	"game-event-service/internal/model"
)

type EventWorker struct {
	events      chan model.Event
	workerCount int
	metrics     *metrics.Metrics
	wg          sync.WaitGroup
}

func NewEventWorker(
	bufferSize int,
	workerCount int,
	metrics *metrics.Metrics,
) *EventWorker {
	return &EventWorker{
		events:      make(chan model.Event, bufferSize),
		workerCount: workerCount,
		metrics:     metrics,
	}
}

func (w *EventWorker) Run(ctx context.Context) {
	for i := 1; i <= w.workerCount; i++ {
		w.wg.Add(1)

		go func(id int) {
			defer w.wg.Done()

			log.Printf("worker %d started", id)

			for {
				select {
				case event := <-w.events:
					log.Printf(
						"worker %d processing event: type=%s player_id=%d",
						id,
						event.Type,
						event.PlayerID,
					)

					w.metrics.RecordEventProcessed()

				case <-ctx.Done():
					log.Printf("worker %d stopped", id)
					return
				}
			}
		}(i)
	}

	w.wg.Wait()
}

func (w *EventWorker) Submit(event model.Event) {
	w.events <- event
}
