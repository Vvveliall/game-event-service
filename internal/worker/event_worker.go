package worker

import (
	"context"
	"log"
	"sync"

	"game-event-service/internal/model"
)

type EventWorker struct {
	events      chan model.Event
	workerCount int
	wg          sync.WaitGroup
}

func NewEventWorker(bufferSize, workerCount int) *EventWorker {
	return &EventWorker{
		events:      make(chan model.Event, bufferSize),
		workerCount: workerCount,
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
