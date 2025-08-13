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

	"github.com/joho/godotenv"
	"oprosdom.ru/msvc_codesender/internal/gateway"
	"oprosdom.ru/msvc_codesender/internal/handlers"
	"oprosdom.ru/msvc_codesender/internal/models"
	"oprosdom.ru/msvc_codesender/internal/repo"
	"oprosdom.ru/msvc_codesender/internal/service"
	http_client "oprosdom.ru/msvc_codesender/internal/transport/http"
)

func init() {
	// в проде не будет .env файла, поэтому делаем такую проверку: на локале берется .env, в проде из k8s
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
}

func main() {

	zimaURL := os.Getenv("ZIMA_URL")
	zimaKey := os.Getenv("ZIMA_KEY")
	zimaDevice := os.Getenv("ZIMA_DEVICE")

	// // Загружаем .env файл
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	cfg := models.Config{
		WorkerInterval: 5 * time.Second,
		MaxJitterMs:    500,
		Gateways: []models.GatewayConfig{
			{
				Name: "Zima1reg",
				URL:  zimaURL,
				Type: "regular",
				Auth: map[string]string{"token": zimaKey, "device": zimaDevice},
			},
			{
				Name: "Zima2reg",
				URL:  zimaURL,
				Type: "regular",
				Auth: map[string]string{"token": zimaKey, "device": zimaDevice},
			},
			{
				Name: "Zima3prem",
				URL:  zimaURL,
				Type: "premium",
				Auth: map[string]string{"token": zimaKey, "device": zimaDevice},
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

	postgresRepo, err := repo.NewRepoFactory(ctx, repoConn)
	if err != nil {
		log.Fatalf("repo initialization failed with error: %v", err)
	}
	defer postgresRepo.Close() // это важно чтоб при закрытии разрывать соединения с базой иначе при многократном рестарте приложения лимит подключений к postgresql иссякнет и получим too many connections

	noSqlRepo, err := repo.NewNoSqlRepoFactory(ctx, "mongodb://admin:admin@localhost:27017", "logs")
	if err != nil {
		log.Fatalf("nosql initialization failed with error: %v", err)
	}
	defer noSqlRepo.Close(ctx)

	transport := http_client.NewHTTPTransport()

	svc := service.NewServiceFactory(postgresRepo)

	// предусмотреть контекст!
	// 1) реализовать graceful shutdown так, чтоб на начатые запросы завершались, а новые не принимались. Чтоб не получилось так, что из кафки возьмем сообщение и завершимся. Нужно обязательно чтоб отправка происходила.

	// в случае с http сервером мы загоняли сервис в хендлер, здесь же мы загоняем сервис в kafkaConsumer
	kafkaConsumer := handlers.NewConsumer([]string{"localhost:9092"}, "code", "for_notify", svc)

	// Горутина для запуска и перезапуска Kafka consumer
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Kafka consumer shutdown requested.")
				return // это выход из горутины
			default:
				log.Println("Starting Kafka consumer...")
				if err := kafkaConsumer.Start(ctx); err != nil {
					log.Printf("Kafka consumer error: %v", err)
					time.Sleep(5 * time.Second) // Пауза перед повтором
					continue
				}
				// В идеале мы сюда вообще не дойдем, но на всякий случай вдруг что при отмене контекста произойдет
				log.Println("Kafka consumer stop work without error")
				return // это выход из горутины
			}
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
			Repo:      noSqlRepo,
			Config:    gwCfg.Auth,
		}

		worker := handlers.NewWorker(
			fmt.Sprintf("worker-%d", i+1),
			gateway,
			cfg.WorkerInterval,
			cfg.MaxJitterMs,
			svc,
		)

		wg.Add(1)
		go worker.Run(ctx, &wg)
	}

	// Канал для контроля завершения всех воркеров
	workersStop := make(chan struct{})

	go func() {
		wg.Wait()
		close(workersStop)
	}()

	// select блокирует программу до наступления одного из case. бесконечный for здесь смысла не имеет, так как select запускается единожды и результат любого case приводит к завершению программы
	// select {
	// case <-ctx.Done():

	<-ctx.Done()
	// Даем разумное время на завершение обработки отправки сообщений
	// ctx контекст не передаем иначе shutdownCtx немедленно отменится при отмене ctx. Нужно делать новый контекст чтоб дать время
	// TODO проработать вопрос с контекстом shutdownCtx ввиду того, что в Transport.Post мы используем независимый контекст, чтоб дооделать взятые задачи
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Начинаем завершение всех процессов...")

	select {
	case <-workersStop: // срабатывает когда канал закрыт
		log.Println("Graceful shutdown done correctly")
	case <-shutdownCtx.Done():
		log.Println("Graceful shutdown done by timeout")
	}
	// }

}
