package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"oprosdom.ru/microservice_notify/internal/gateway"
	"oprosdom.ru/microservice_notify/internal/models"
	"oprosdom.ru/microservice_notify/internal/repo"
	"oprosdom.ru/microservice_notify/internal/service"
	http_client "oprosdom.ru/microservice_notify/internal/transport/http"
	"oprosdom.ru/microservice_notify/internal/transport/kafka"
	"oprosdom.ru/microservice_notify/internal/worker"
)

func main() {

	cfg := models.Config{
		WorkerInterval: 5 * time.Second,
		MaxJitterMs:    500,
		GatewayTimeout: 30 * time.Second,
		Gateways: []models.GatewayConfig{
			{
				Name: "API1",
				URL:  "https://api1.example.com/send",
				Type: "regular",
				Auth: map[string]string{"api_key": "api1_key"},
			},
			{
				Name: "API2",
				URL:  "https://api2.example.com/sms",
				Type: "regular",
				Auth: map[string]string{"token": "api2_token"},
			},
			{
				Name: "API3",
				URL:  "https://api3.example.com/message",
				Type: "premium",
				Auth: map[string]string{"auth_key": "api3_auth"},
			},
		},
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repoConn := "postgres://test:test@127.0.0.1:5432/notify?" +
		"sslmode=disable&" +
		"pool_min_conns=5&" +
		"pool_max_conns=25&" +
		"pool_max_conn_lifetime=30m&" +
		"pool_max_conn_lifetime_jitter=5m&" +
		"pool_max_conn_idle_time=15m&" +
		"pool_health_check_period=1m"

	repo, err := repo.NewRepoFactory(ctx, repoConn)
	if err != nil {
		log.Fatalf("repo initialization failed with error: %v", err)
	}
	defer repo.Close() // это важно чтоб при закрытии разрывать соединения с базой иначе при многократном рестарте приложения лимит подключений к postgresql иссякнет и получим too many connections

	transport := http_client.NewHTTPTransport(cfg.GatewayTimeout)

	// codeTransport, err := transport.NewTransportFactory(ctx, "localhost:9092", "code")
	// if err != nil {
	// 	log.Fatalf("codeTransport initialization failed with error: %v", err)
	// }
	// defer codeTransport.Close()

	svc := service.NewServiceFactory(repo)

	// предусмотреть контекст!
	// 1) реализовать graceful shutdown так, чтоб на начатые запросы завершались, а новые не принимались. Чтоб не получилось так, что из кафки возьмем сообщение и завершимся. Нужно обязательно чтоб отправка происходила.

	// в случае с http сервером мы загоняли сервис в хендлер, здесь же мы загоняем сервис в kafkaConsumer
	//kafkaConsumer := kafka.NewConsumer([]string{"localhost:9092"}, cfg.KafkaTopic, cfg.KafkaGroupID, svc)
	kafkaConsumer := kafka.NewConsumer([]string{"localhost:9092"}, "code", "for_notify", svc)

	errCh := make(chan error, 1)

	go func() {
		if err := kafkaConsumer.Start(ctx); err != nil {
			errCh <- err
		}
	}()

	var wg sync.WaitGroup

	// Создаем по одному воркеру для каждого шлюза
	for i, gwCfg := range cfg.Gateways {
		gateway := &gateway.Gateway{
			Name:      gwCfg.Name,
			URL:       gwCfg.URL,
			Type:      gwCfg.Type,
			Transport: transport,
			Config:    gwCfg.Auth,
		}

		worker := worker.NewWorker(
			fmt.Sprintf("worker-%d", i+1),
			gateway,
			repo,
			cfg.WorkerInterval,
			cfg.MaxJitterMs,
		)

		wg.Add(1)
		go worker.Run(ctx, &wg)
	}

	// select блокирует программу до наступления одного из case. бесконечный for здесь смысла не имеет, так как select запускается единожды и результат любого case приводит к завершению программы
	select {
	case err := <-errCh: // присваиваем err значение из канала, область видимости только этот case
		log.Fatalf("Возникла ошибка: %v", err)
	case <-ctx.Done():
		// Даем разумное время на завершение обработки отправки sms
		// ctx контекст не передаем иначе shutdownCtx немедленно отменится при отмене ctx. Нужно делать новый контекст чтоб дать время
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Println("Начинаем завершение всех процессов...")

		// Канал для контроля завершения всех воркеров
		done := make(chan struct{})

		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done: // срабатывает когда канал закрыт
			log.Println("Graceful shutdown done correctly")
		case <-shutdownCtx.Done():
			log.Println("Graceful shutdown done by timeout")
		}
	}

}
