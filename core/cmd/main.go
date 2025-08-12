// @title           Core API
// @version         1.0
// @description     Аутентификация и выдача токенов осуществляется /auth
// @host            localhost:8082
// @BasePath        /api/v1
// @schemes         http
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name auth

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

	httpSwagger "github.com/swaggo/http-swagger"
	_ "oprosdom.ru/swagger/core"

	polls_handlers "oprosdom.ru/core/internal/polls/handlers"
	polls_repo "oprosdom.ru/core/internal/polls/repo"
	polls_service "oprosdom.ru/core/internal/polls/service"
	users_handlers "oprosdom.ru/core/internal/users/handlers"
	users_repo "oprosdom.ru/core/internal/users/repo"
	users_service "oprosdom.ru/core/internal/users/service"
)

var publicKey *rsa.PublicKey

// TODO
// func init() {
//     // загрузка ключа при старте
//     publicKey = initPublicKey()
// }

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

	//////////////////////////////////////// POLLS INITIALIZATION ////////////////////////////////////////

	pollsRepoConn := "postgres://test:test@127.0.0.1:5436/polls?" +
		"sslmode=disable&" +
		"pool_min_conns=5&" +
		"pool_max_conns=25&" +
		"pool_max_conn_lifetime=30m&" +
		"pool_max_conn_lifetime_jitter=5m&" +
		"pool_max_conn_idle_time=15m&" +
		"pool_health_check_period=1m"

	pollsRepo, err := polls_repo.NewRepoFactory(ctx, pollsRepoConn)
	if err != nil {
		log.Fatalf("polls repo initialization failed with error: %v", err)
	}
	defer pollsRepo.Close() // это важно чтоб при закрытии разрывать соединения с базой иначе при многократном рестарте приложения лимит подключений к postgresql иссякнет и получим too many connections

	pollsService := polls_service.NewServiceFactory(pollsRepo)
	pollsHandler := polls_handlers.NewHandler(pollsService)

	//////////////////////////////////////////////////////////////////////////////////////////////////////

	// предусмотреть контекст!
	// 1) если клиент стопнул в браузере выполнение, то нужно отменять операции -> это предусмотрено http сервером, но нужно обрабатывать  это событие в хендлерах
	// 2) реализовать graceful shutdown так, чтоб на начатые запросы завершались, а новые не принимались

	// curl -X POST "http://127.0.0.1:8082/api/poll/vote" -H "Content-Type: application/json" -d '{"poll_id":1,"vote":"za"}' -b cookies.txt
	// curl -X POST "http://127.0.0.1:8080/api/kvartira" -H "Content-Type: application/json" -d '{"number":"115","komnat":2}'
	// curl -X PUT "http://127.0.0.1:8080/api/member/09770685-ae8a-4e68-9751-50bba1d846f1" -H "Content-Type: application/json" -d '{"name":"Иван Иванов","phone":"+79991234567","community":5}'
	// curl -X PUT "http://127.0.0.1:8080/api/kvartira/1f970b7b-679c-4c7d-a252-3ef370d439f4" -H "Content-Type: application/json" -d '{"number":"12","komnat":3}'
	// curl -X GET "http://127.0.0.1:8082/api/v1/polls/stat" -H "Content-Type: application/json" -b cookies.txt

	r := mux.NewRouter()

	requireJwt := r.PathPrefix("/api/v1").Subrouter()
	requireJwt.Use(jwtMiddleware) // это применение промежуточного хендлера до вызова основного. Можем добавлять промеж хендлеры цепочкой

	requireJwt.Use(jwtMiddleware)
	requireJwt.HandleFunc("/users/verification/phone", usersHandler.PhoneSend).Methods("POST")
	requireJwt.HandleFunc("/polls", pollsHandler.GetPolls).Methods("GET")
	requireJwt.HandleFunc("/polls/stat", pollsHandler.PollStats).Methods("GET")
	requireJwt.HandleFunc("/poll/vote", pollsHandler.Vote).Methods("POST")

	// Swagger UI без middleware
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

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

		// вообще считается антипаттерном, но пока не нашел способ как из мидлвары передать в основной хендлер

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if jti, exists := claims["jti"].(string); exists {
				ctx := context.WithValue(r.Context(), "jti", jti)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}
