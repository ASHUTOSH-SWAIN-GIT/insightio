package ingest

import (
	"log"

	"github.com/ASHUTOSH-SWAIN-GIT/insightio/internal/metrics/store"

	pb "github.com/ASHUTOSH-SWAIN-GIT/insightio/proto"
)

// Worker pulls full Event objects from eventChan and updates the MetricStore.
type Worker struct {
	eventChan   <-chan *pb.Event   // receives events
	metricStore *store.MetricStore // reference to metric store
	stopChan    chan struct{}      // for graceful shutdown
}

// NewWorker creates a worker bound to eventChan and metric store.
func NewWorker(eventChan <-chan *pb.Event, store *store.MetricStore) *Worker {
	return &Worker{
		eventChan:   eventChan,
		metricStore: store,
		stopChan:    make(chan struct{}),
	}
}

// Start launches the worker loop in a goroutine.
func (w *Worker) Start() {
	go func() {
		log.Println("Worker started")

		for {
			select {

			// event received from ingest service
			case event, ok := <-w.eventChan:
				if !ok {
					log.Println("eventChan closed, worker stopping")
					return
				}
				if event != nil {
					w.metricStore.AddEvent(event.Type)
				}

			// stop signal received
			case <-w.stopChan:
				log.Println("Worker stopped")
				return
			}
		}
	}()
}

// Stop gracefully shuts down the worker.
func (w *Worker) Stop() {
	close(w.stopChan)
}
