// @title           AUTH API
// @version         1.0
// @description     Аутентификация и выдача токенов
// @host            localhost:8081
// @BasePath        /auth
// @schemes         http

package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "oprosdom.ru/swagger/auth"

	"github.com/gorilla/mux"
	"oprosdom.ru/msvc_auth/internal/handlers"
	"oprosdom.ru/msvc_auth/internal/models"
	"oprosdom.ru/msvc_auth/internal/repo"
	"oprosdom.ru/msvc_auth/internal/service"
	"oprosdom.ru/msvc_auth/internal/transport"
	"oprosdom.ru/shared"
	"oprosdom.ru/shared/models/pb/access"
)

var healthOK = []byte("OK\n") // заранее подготовленный ответ для health проверки k8s

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	authDbURI := os.Getenv("AUTH_DB_URI")
	redisAddr := os.Getenv("REDIS_URI")
	kafkaURI := os.Getenv("KAFKA_URI")
	msvcAccessURI := os.Getenv("MSVC_ACCESS_URI")

	keyPath := os.Getenv("PRIVATE_KEY_PATH")
	if keyPath == "" {
		keyPath = "private.pem" // fallback для локальной разработки
	}

	privateKey, err := loadPrivateKey(keyPath)

	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	pubkeyId := shared.GetPubkeyId(&privateKey.PublicKey)

	key := &models.KeyData{
		PrivateKey: privateKey,
		PubkeyId:   pubkeyId,
	}

	log.Printf("KeyID: %s", pubkeyId)

	postgresConn := authDbURI +
		"&pool_min_conns=5&" +
		"pool_max_conns=25&" +
		"pool_max_conn_lifetime=30m&" +
		"pool_max_conn_lifetime_jitter=5m&" +
		"pool_max_conn_idle_time=15m&" +
		"pool_health_check_period=1m"

	// на будущее! чтоб не забыть! нельзя называть переменную также как пакет, иначе еще раз этот пакет не вызвать!
	postgres, err := repo.NewRepoFactory(ctx, postgresConn)
	if err != nil {
		log.Fatalf("postgresql initialization failed with error: %v", err)
	}
	defer postgres.Close() // это важно чтоб при закрытии разрывать соединения с базой иначе при многократном рестарте приложения лимит подключений к postgresql иссякнет и получим too many connections

	redis, err := repo.NewRamRepoFactory(ctx, redisAddr)
	if err != nil {
		log.Fatalf("redis initialization failed with error: %v", err)
	}
	defer redis.Close()

	codeTransport, err := transport.NewTransportFactory(ctx, kafkaURI, "code")
	if err != nil {
		log.Fatalf("codeTransport initialization failed with error: %v", err)
	}
	defer codeTransport.Close()

	grpcClient := transport.NewGrpcClient(msvcAccessURI)
	if err := grpcClient.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer grpcClient.Close()

	// в сервис передаем NewAccessClient это автоматически сгенерированный конструктор, который возвращает объект с методами, которые будут исполняться на сервере grpc
	accessClient := access.NewAccessClient(grpcClient.Connection())

	authService := service.NewServiceFactory(key, redis, postgres, codeTransport, accessClient)
	h := handlers.NewHandler(authService)

	// предусмотреть контекст!
	// 1) если клиент стопнул в браузере выполнение, то нужно отменять операции -> это предусмотрено http сервером, но нужно обрабатывать  это событие в хендлерах (по сути перед затратными операциями нужно ловить отмену контекста)
	// 2) реализовать graceful shutdown так, чтоб на начатые запросы завершались, а новые не принимались

	// curl -k -X POST "https://127.0.0.1/auth/phone" -H "Content-Type: application/json" -d '{"phone":"+79994951548"}'
	// curl -c cookies.txt -k -X POST "https://127.0.0.1/auth/code" -H "Content-Type: application/json" -d '{"phone":"+79994951548", "code":1234}'

	r := mux.NewRouter()
	r.HandleFunc("/auth/phone", h.PhoneSend).Methods("POST")
	r.HandleFunc("/auth/code", h.CodeCheck).Methods("POST")
	// k8s health check
	r.HandleFunc("/health", healthHandler).Methods("GET")
	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
		// таймауты read/write тут не указываем потому что у нас вероятно будут вебсокеты на этом же порту, поэтому будем другими механизмами таймауты отслеживать
		// хотя нужно изучить вопрос, ws это же не http, а поверх http
	}

	errCh := make(chan error, 1)

	// запускаем в отдельной горутине потому что ListenAndServe это блокирующий вызов, и без горутины мы до select никогдай не дойдем. А значит graceful shutdown не получится.
	go func() {
		log.Printf("Стартуем noTLS микросервис AUTH на %v порту", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	// select блокирует программу до наступления одного из case. бесконечный for здесь смысла не имеет, так как select запускается единожды и результат любого case приводит к завершению программы
	select {
	case err := <-errCh: // присваиваем err значение из канала, область видимости только этот case
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Println("Сервак остановлен. Это не неожиданная ошибка, а ожидаемое поведение завершения")
		default:
			log.Fatalf("Возникла ошибка: %v", err)
		}
	case <-ctx.Done():
		// мы не должны передавать ctx иначе shutdownCtx немедленно отменится при отмене ctx. Вместо этого делаем новый контекст, который даст 5 сек на довыполнение
		shutdownCtx, shutdownCtxStop := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCtxStop()

		// Shutdown() перестанет принимать новые подключения, но завершит старые
		// на заметку Shutdown() не ждёт завершения WS-соединений
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful Shutdown http сервера не сработало, ошибка: %v", err)
		}
		log.Println("Graceful Shutdown http сервера успешен")
	}

}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	pemData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaKey, nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(healthOK)
}
