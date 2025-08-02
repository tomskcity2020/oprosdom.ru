package transport

import (
	"context"

	kafka "oprosdom.ru/microservice_auth/internal/transport/internal"
	shared_models "oprosdom.ru/shared/models"
)

func NewTransportFactory(ctx context.Context, conn string, topic string) (TransportInterface, error) {
	return kafka.NewKafka(ctx, conn, topic)
}

type TransportInterface interface {
	Send(ctx context.Context, msg *shared_models.MsgCode) error
	Close() error
}
