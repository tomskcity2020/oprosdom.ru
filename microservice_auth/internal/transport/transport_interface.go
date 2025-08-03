package transport

import (
	"context"

	kafka "oprosdom.ru/microservice_auth/internal/transport/internal"
	pb "oprosdom.ru/shared/models/proto"
)

func NewTransportFactory(ctx context.Context, conn string, topic string) (TransportInterface, error) {
	return kafka.NewKafka(ctx, conn, topic)
}

type TransportInterface interface {
	Send(ctx context.Context, msg *pb.MsgCode) error
	Close() error
}
