package transport

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	conn   *grpc.ClientConn
	mu     sync.RWMutex
	target string
}

func NewGrpcClient(target string) *GrpcClient {
	return &GrpcClient{
		target: target,
	}
}

func (c *GrpcClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// чекаем существующее соединение и его состояние
	if c.conn != nil {
		state := c.conn.GetState()
		if state == connectivity.Ready || state == connectivity.Connecting || state == connectivity.Idle {
			return nil
		}
		// если дошли сюда, значит с соединение что-то не то - закрываем его
		c.conn.Close()
		c.conn = nil
	}

	// TODO нужно заменить insecure на tls
	conn, err := grpc.NewClient(
		c.target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(dialCtx context.Context, addr string) (net.Conn, error) {
			// используем переданный контекст с таймаутом
			return (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext(dialCtx, "tcp", addr)
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", c.target, err)
	}

	// соединение устанавливается при первом RPC вызове поэтому здесь просто сохраняем соединение

	c.conn = conn
	return nil
}

func (c *GrpcClient) Connection() *grpc.ClientConn {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// чекаем состояние
	if c.conn != nil {
		state := c.conn.GetState()
		if state == connectivity.Shutdown || state == connectivity.TransientFailure {
			return nil
		}
	}

	return c.conn
}

func (c *GrpcClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// значит уже закрыто
	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.conn = nil
	return err
}
