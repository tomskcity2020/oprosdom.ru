package main

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"

	// "oprosdom.ru/core/internal/dz/internal/handlers"
	users_handlers "oprosdom.ru/core/internal/users/handlers"
	users_repo "oprosdom.ru/core/internal/users/repo"
	users_service "oprosdom.ru/core/internal/users/service"
)

var publicKey *rsa.PublicKey

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := initPublicKey(); err != nil {
		log.Fatalf("Failed to initialize public key: %v", err)
	}

	//////////////////////////////////////// USERS INITIALIZATION ////////////////////////////////////////

	userRepoConn := "postgres://test:test@127.0.0.1:5435/users?" +
		"sslmode=disable&" +
		"pool_min_conns=5&" +
		"pool_max_conns=25&" +
		"pool_max_conn_lifetime=30m&" +
		"pool_max_conn_lifetime_jitter=5m&" +
		"pool_max_conn_idle_time=15m&" +
		"pool_health_check_period=1m"

	usersRepo, err := users_repo.NewRepoFactory(ctx, userRepoConn)
	if err != nil {
		log.Fatalf("user repo initialization failed with error: %v", err)
	}
	defer usersRepo.Close() // это важно чтоб при закрытии разрывать соединения с базой иначе при многократном рестарте приложения лимит подключений к postgresql иссякнет и получим too many connections

	usersService := users_service.NewServiceFactory(usersRepo)
	usersHandler := users_handlers.NewHandler(usersService)

	//////////////////////////////////////////////////////////////////////////////////////////////////////

	// предусмотреть контекст!
	// 1) если клиент стопнул в браузере выполнение, то нужно отменять операции -> это предусмотрено http сервером, но нужно обрабатывать  это событие в хендлерах
	// 2) реализовать graceful shutdown так, чтоб на начатые запросы завершались, а новые не принимались

	// curl -X POST "http://127.0.0.1:8080/api/member" -H "Content-Type: application/json" -d '{"name":"Иван Иванов","phone":"+79991234567","community":5}'
	// curl -X POST "http://127.0.0.1:8080/api/kvartira" -H "Content-Type: application/json" -d '{"number":"115","komnat":2}'
	// curl -X PUT "http://127.0.0.1:8080/api/member/09770685-ae8a-4e68-9751-50bba1d846f1" -H "Content-Type: application/json" -d '{"name":"Иван Иванов","phone":"+79991234567","community":5}'
	// curl -X PUT "http://127.0.0.1:8080/api/kvartira/1f970b7b-679c-4c7d-a252-3ef370d439f4" -H "Content-Type: application/json" -d '{"number":"12","komnat":3}'
	// curl -X GET "http://127.0.0.1:8080/api/members" -H "Content-Type: application/json"
	// curl -X GET "http://127.0.0.1:8080/api/kvartiras" -H "Content-Type: application/json"
	// curl -X GET "http://127.0.0.1:8080/api/member/09770685-ae8a-4e68-9751-50bba1d846f1"
	// curl -X GET "http://127.0.0.1:8080/api/kvartira/1f970b7b-679c-4c7d-a252-3ef370d439f4"
	// curl -X DELETE "http://127.0.0.1:8080/api/member/09770685-ae8a-4e68-9751-50bba1d846f1"
	// curl -X DELETE "http://127.0.0.1:8080/api/kvartira/1f970b7b-679c-4c7d-a252-3ef370d439f4"
	// curl -X POST "http://127.0.0.1:8080/api/member/b406690d-018a-4fb5-b1b7-70946b92432e/paydebt" -H "Content-Type: application/json" -d '{"kvartira_id":"f983de63-e4f3-483a-b9a2-6e1f3ea960b7","amount":"17500.55"}'

	r := mux.NewRouter()
	r.Use(jwtMiddleware) // это применение промежуточного хендлера до вызова основного. Можем добавлять промеж хендлеры цепочкой
	r.HandleFunc("/api/users/verification/phone", usersHandler.PhoneSend).Methods("POST")
	// verification/phonecode
	// verification/bankcard
	// r.HandleFunc("/api/kvartira", h.KvartiraAdd).Methods("POST")
	// r.HandleFunc("/api/member/{id}", h.MemberUpdate).Methods("PUT")
	// r.HandleFunc("/api/kvartira/{id}", h.KvartiraUpdate).Methods("PUT")
	// r.HandleFunc("/api/members", h.MembersGet).Methods("GET")
	// r.HandleFunc("/api/kvartiras", h.KvartirasGet).Methods("GET")
	// r.HandleFunc("/api/member/{id}", h.MemberGet).Methods("GET")
	// r.HandleFunc("/api/kvartira/{id}", h.KvartiraGet).Methods("GET")
	// r.HandleFunc("/api/{mk}/{id}", h.RemoveById).Methods("DELETE")
	// r.HandleFunc("/api/member/{id}/paydebt", h.PayDebt).Methods("POST")

	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
		// таймауты read/write тут не указываем потому что у нас вероятно будут вебсокеты на этом же порту, поэтому будем другими механизмами таймауты отслеживать
		// хотя нужно изучить вопрос, ws это же не http, а поверх http
	}

	errCh := make(chan error, 1)

	// запускаем в отдельной горутине потому что ListenAndServe это блокирующий вызов, и без горутины мы до select никогдай не дойдем. А значит graceful shutdown не получится.
	go func() {
		log.Printf("Стартуем noTLS сервак на %v порту", srv.Addr)
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

func initPublicKey() error {
	pubKeyPEM := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwNoGHw2jZBJp77npZDc3
uxM/uzEq8Gd4myO9+wQG8VeNlITOoXWBx77LRe2TXZQYyRc9jYudW3xH2WvZ49m2
w8hnmCfwRA+tOkdzoosLnQ1dxom+kVbIXiDIefmmnoXhvnXEg9jAeB+csSnmyD+Q
vsV1/U13/O7iRcnUxU3mkF5knEPwQ1GXUH9Aiv5YJC2JcEegsAq2hLCosd1eCytg
A/9FMmyA7qWAQHlX3jai2p91SOtO6OROEmZ3MeaxTv0T4vforyqy+cPNaCS6GC3U
Zt8jqQg7y5HEXpvvb60fp3hYge5FDa8Cug0wzMgitUTxZ8y4VtUofPuHCjB7YgBS
+QIDAQAB
-----END PUBLIC KEY-----`
	// pem.Decode отдельно делать не нужно, потому что под капотом ParseRSAPublicKeyFromPEM тот же pem.decode
	pub, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKeyPEM))
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	publicKey = pub
	log.Println("public key init success")
	return nil
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil {
			// описание не возвращаем, в этом нет смысла в нашем случае
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			// если приведение к типу не сработает, то будет ошибка. Обязательная проверка иначе могут в alg токена прислать что угодно
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected jwt signing alg")
			}
			return publicKey, nil
		})

		// если токен не прошел проверку подписи или истек - token.Valud будет false
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
