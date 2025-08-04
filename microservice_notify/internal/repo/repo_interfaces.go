package repo

import (
	"context"

	"oprosdom.ru/microservice_notify/internal/models"
	repo_internal "oprosdom.ru/microservice_notify/internal/repo/internal"
)

func NewRepoFactory(ctx context.Context, conn string) (RepositoryInterface, error) {
	return repo_internal.NewPostgres(ctx, conn)
}

type RepositoryInterface interface {
	Close()
	CreateMessage(ctx context.Context, msg models.SMSMessage) error
	GetNextMessageForGateway(ctx context.Context, gatewayType string) (*models.SMSMessage, error)
	UpdateMessageStatus(ctx context.Context, id int, status, gateway string) error
	LogAttempt(ctx context.Context, msgID int, workerID, gateway, status, errorMsg string) error
}
