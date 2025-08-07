package service_internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"oprosdom.ru/msvc_codesender/internal/gateway"
)

func (s *ServiceStruct) ProcessMessage(ctx context.Context, gateway *gateway.Gateway) error {

	// не берем новые месаги если контекст отменен
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	// если передавать ctx в repo, то мы нарушим тогда логику программы: если придет отмена контекста, то мы перестанем брать новые дела, но этот же ctx уйдет в репо и сорвет там уже начатую операцию и тот же update отменится
	waitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg, err := s.repo.GetNextSmsForGateway(waitCtx, gateway.Type)
	if err != nil {
		return fmt.Errorf("worker error: %w", err)
	}
	if msg == nil {
		return errors.New("empty msg")
	}

	s.repo.LogSmsAttempt(waitCtx, msg.Id, gateway.Name, "attempt", "")

	if err = gateway.Send(waitCtx, *msg); err != nil {
		s.repo.LogSmsAttempt(waitCtx, msg.Id, gateway.Name, "failed", err.Error())

		if updateErr := s.repo.UpdateMessageStatus(waitCtx, msg.Id, "wait", ""); updateErr != nil {
			log.Printf("Worker failed to return wait status: %v", updateErr)
		}
		return fmt.Errorf("worker send failed: %w", err)
	}

	if err := s.repo.UpdateMessageStatus(waitCtx, msg.Id, "sent", gateway.Name); err != nil {
		log.Printf("Worker failed to set sent status: %v", err)
	}
	s.repo.LogSmsAttempt(waitCtx, msg.Id, gateway.Name, "sent", "")

	return nil

}
