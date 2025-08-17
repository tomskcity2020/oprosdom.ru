package transport

import (
	"context"

	"google.golang.org/grpc"
	transport "oprosdom.ru/msvc_auth/internal/transport/internal"
	"oprosdom.ru/shared/models/pb"
)

func NewTransportFactory(ctx context.Context, conn string, topic string) (TransportInterface, error) {
	return transport.NewKafka(ctx, conn, topic)
}

type TransportInterface interface {
	Send(ctx context.Context, msg *pb.MsgCode) error
	Close() error
}

func NewGrpcClient(target string) GrpcClientInterface {
	return transport.NewGrpcClient(target)
}

type GrpcClientInterface interface {
	Connect(ctx context.Context) error
	Connection() *grpc.ClientConn
	Close() error
}
