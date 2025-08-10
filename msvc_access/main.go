package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"oprosdom.ru/msvc_access/internal/repo"
	"oprosdom.ru/msvc_access/internal/service"
	"oprosdom.ru/shared/models/pb/access"
)

// адаптер для grpc сервера, делегирование вызовов сервисному слою
type accessServer struct {
	access.UnimplementedAccessServer
	service service.Service
}

func (s *accessServer) AddWhitelist(ctx context.Context, req *access.SendRequest) (*access.SendResponse, error) {
	return s.service.AddWhitelist(ctx, req)
}

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	redisAddr := "localhost:6380"
	redis, err := repo.NewRamRepoFactory(ctx, redisAddr)
	if err != nil {
		log.Fatalf("redis initialization failed with error: %v", err)
	}
	defer redis.Close()

	accessService := service.NewServiceFactory(redis)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	access.RegisterAccessServer(grpcServer, &accessServer{
		service: accessService,
	})

	errCh := make(chan error, 1)

	// запускаем в отдельной горутине потому что Serve() блокирующий вызов, и без горутины мы до select никогдай не дойдем. А значит graceful shutdown не получится.
	go func() {

		log.Printf("Access server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			errCh <- err
		}

	}()

	// select блокирует программу до наступления одного из case. бесконечный for здесь смысла не имеет, так как select запускается единожды и результат любого case приводит к завершению программы
	select {
	case err := <-errCh: // присваиваем err значение из канала, область видимости только этот case
		log.Fatalf("Возникла ошибка: %v", err)
	case <-ctx.Done():
		log.Println("Graceful Shutdown started")
		grpcServer.GracefulStop()
		log.Println("GRPC server stopped")
	}

}
