package worker

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"oprosdom.ru/microservice_notify/internal/gateway"
	"oprosdom.ru/microservice_notify/internal/repo"
)

type Worker struct {
	ID          string
	Gateway     *gateway.Gateway
	Repo        repo.RepositoryInterface
	Interval    time.Duration
	rand        *rand.Rand
	maxJitterMs int
}

func NewWorker(id string, gateway *gateway.Gateway, repo repo.RepositoryInterface, interval time.Duration, maxJitterMs int) *Worker {
	src := rand.NewSource(time.Now().UnixNano())
	return &Worker{
		ID:          id,
		Gateway:     gateway,
		Repo:        repo,
		Interval:    interval,
		rand:        rand.New(src),
		maxJitterMs: maxJitterMs,
	}
}

func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	initialJitter := time.Duration(w.rand.Intn(w.maxJitterMs)) * time.Millisecond
	select {
	case <-time.After(initialJitter):
	case <-ctx.Done():
		return
	}

	log.Printf("Worker %s started for gateway %s (%s)", w.ID, w.Gateway.Name, w.Gateway.Type)

	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %s stopping", w.ID)
			return
		case <-ticker.C:
			w.processMessage(ctx)

			jitter := time.Duration(w.rand.Intn(w.maxJitterMs)) * time.Millisecond
			select {
			case <-time.After(jitter):
			case <-ctx.Done():
				return
			}
		}
	}
}

func (w *Worker) processMessage(ctx context.Context) {
	// Используем тип шлюза вместо его имени
	msg, err := w.Repo.GetNextMessageForGateway(ctx, w.Gateway.Type)
	if err != nil {
		log.Printf("Worker %s error: %v", w.ID, err)
		return
	}
	if msg == nil {
		return
	}

	log.Printf("Worker %s processing message ID %d", w.ID, msg.ID)

	w.Repo.LogAttempt(ctx, msg.ID, w.ID, w.Gateway.Name, "attempt", "")

	err = w.Gateway.Send(ctx, *msg)
	if err == nil {
		if err := w.Repo.UpdateMessageStatus(ctx, msg.ID, "sent", w.Gateway.Name); err != nil {
			log.Printf("Worker %s update status error: %v", w.ID, err)
		}
		w.Repo.LogAttempt(ctx, msg.ID, w.ID, w.Gateway.Name, "sent", "")
		return
	}

	w.Repo.LogAttempt(ctx, msg.ID, w.ID, w.Gateway.Name, "failed", err.Error())
	log.Printf("Worker %s failed to send: %v", w.ID, err)

	if updateErr := w.Repo.UpdateMessageStatus(ctx, msg.ID, "not_sent", ""); updateErr != nil {
		log.Printf("Worker %s update status error: %v", w.ID, updateErr)
	}
}
