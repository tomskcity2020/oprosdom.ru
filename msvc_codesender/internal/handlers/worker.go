package handlers

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"oprosdom.ru/msvc_codesender/internal/gateway"
	"oprosdom.ru/msvc_codesender/internal/service"
)

type Worker struct {
	ID          string
	Gateway     *gateway.Gateway
	Interval    time.Duration
	maxJitterMs int
	svc         service.ServiceInterface
	rand        *rand.Rand
}

func NewWorker(id string, gateway *gateway.Gateway, interval time.Duration, maxJitterMs int, svc service.ServiceInterface) *Worker {
	src := rand.NewSource(time.Now().UnixNano())
	return &Worker{
		ID:          id,
		Gateway:     gateway,
		Interval:    interval,
		maxJitterMs: maxJitterMs,
		svc:         svc,
		rand:        rand.New(src),
	}
}

func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// стартовать нужно не всех воркеров одновременно, а с разбежкой небольшой, далее у всех таймер одинаков и разбежка сохранится
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
			w.svc.ProcessMessage(ctx, w.Gateway)
		}
	}
}
