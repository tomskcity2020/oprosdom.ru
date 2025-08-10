package transport

import (
	"context"

	kafka "oprosdom.ru/msvc_auth/internal/transport/internal"
	"oprosdom.ru/shared/models/pb"
)

func NewTransportFactory(ctx context.Context, conn string, topic string) (TransportInterface, error) {
	return kafka.NewKafka(ctx, conn, topic)
}

type TransportInterface interface {
	Send(ctx context.Context, msg *pb.MsgCode) error
	Close() error
}
