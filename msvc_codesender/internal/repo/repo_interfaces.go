package repo

import (
	"context"

	"oprosdom.ru/msvc_codesender/internal/models"
	repo_internal "oprosdom.ru/msvc_codesender/internal/repo/internal"
)

func NewRepoFactory(ctx context.Context, conn string) (RepositoryInterface, error) {
	return repo_internal.NewPostgres(ctx, conn)
}

type RepositoryInterface interface {
	Close()
	InsertSms(ctx context.Context, msg *models.ValidatedMsg) error
	InsertCall(ctx context.Context, msg *models.ValidatedMsg) error
	GetNextSmsForGateway(ctx context.Context, gatewayType string) (*models.MsgFromRepo, error)
	UpdateMessageStatus(ctx context.Context, id int, status, gateway string) error
	LogSmsAttempt(ctx context.Context, msgId int, gateway, status, errorMsg string)
}
