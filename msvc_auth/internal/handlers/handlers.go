package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"oprosdom.ru/msvc_auth/internal/models"
	"oprosdom.ru/msvc_auth/internal/service"
	"oprosdom.ru/shared"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func replyError(w http.ResponseWriter, err error, publicErr string, statusCode int) {
	// оригинальный err.Error() логируем в redis с ttl
	// TODO
	log.Println(err.Error())
	w.WriteHeader(statusCode)
	// отправлять error поле в json'e нет смысла, так как ошибка будет по status code определяться
	json.NewEncoder(w).Encode(publicErr)
}

// PhoneSend godoc
// @Summary      Отправка SMS с кодом верификации на телефон
// @Description  Принимает JSON с номером телефона (в формате E.164), UserAgent и IP заполняются сервером для лимитирования.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        phoneSendReq  body      models.UnsafePhoneSendReq  true  "Данные для отправки кода на телефон. Только поле phone приходит от клиента"
// @Success      201  {string}  string  "success"
// @Failure      400  {string}  string  "incorrect_request или incorrect_phonesend при ошибках валидации"
// @Failure      500  {string}  string  "phonesend_wrong_request при ошибках сервера"
// @Router       /phone [post]
func (h *Handler) PhoneSend(w http.ResponseWriter, r *http.Request) {

	// 1) парсим json
	// 2) проводим санитизацию и первичную валидацию
	// 3) формируем структуру и отдаем в сервис
	// 4) ошибки возвращаем: err логируем, publicErr отдаем фронтенду. publicErr не должен содержать чувствительных данных.

	w.Header().Set("Content-Type", "application/json")

	// парсим и наполняем unsafe структуру данными
	var unsafePhoneSendReq models.UnsafePhoneSendReq
	if err := json.NewDecoder(r.Body).Decode(&unsafePhoneSendReq); err != nil {
		replyError(w, err, "incorrect_request", http.StatusBadRequest)
		return
	}

	unsafePhoneSendReq.UserAgent = r.UserAgent()
	unsafePhoneSendReq.Ip = shared.IpHttpGet(r)

	// проводим первичную валидацию, на выходе заполняем valid структуру и отдаем в сервисный слой (для дальнейшей бизнес-валидации, бизнес-логики и взаимодействия с репо)
	validatedPhoneSendReq, err := unsafePhoneSendReq.Validate()
	if err != nil {
		replyError(w, err, "incorrect_phonesend", http.StatusBadRequest)
		return
	}

	// service (там репо и бизнес)
	if err := h.service.PhoneSend(r.Context(), validatedPhoneSendReq); err != nil {
		replyError(w, err, "phonesend_wrong_request", http.StatusInternalServerError)
		return
	}

	// успех
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("success")

}

// CodeCheck godoc
// @Summary      Проверка кода верификации и выдача JWT токена в куке
// @Description  Принимает JSON с телефоном (E.164) и кодом (uint32), UserAgent и IP сервер подставляет для логирования и валидации. При успешной проверке в ответе устанавливается HttpOnly cookie "auth" с JWT токеном.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        codeCheckReq  body      models.UnsafeCodeCheckReq  true  "Данные для проверки кода: телефон и код"
// @Success      200  {object}  map[string]string  "{"status": "Token issued in cookie"}"
// @Failure      400  {string}  string  "incorrect_request или codecheck_validation_failed при ошибках валидации"
// @Failure      500  {string}  string  "codecheck_wrong_request или create_token_failed при ошибках сервера"
// @Router       /code [post]
func (h *Handler) CodeCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var unsafeCodeCheckReq models.UnsafeCodeCheckReq
	if err := json.NewDecoder(r.Body).Decode(&unsafeCodeCheckReq); err != nil {
		replyError(w, err, "incorrect_request", http.StatusBadRequest)
		return
	}

	unsafeCodeCheckReq.UserAgent = r.UserAgent()
	unsafeCodeCheckReq.Ip = shared.IpHttpGet(r)

	validatedCodeCheckReq, err := unsafeCodeCheckReq.Validate()
	if err != nil {
		replyError(w, err, "codecheck_validation_failed", http.StatusBadRequest)
		return
	}

	if err := h.service.CodeCheck(r.Context(), validatedCodeCheckReq); err != nil {
		replyError(w, err, "codecheck_wrong_request", http.StatusInternalServerError)
		return
	}

	// 10 лет в часах (для jwt и куки)
	tenYears := time.Hour * 24 * 365 * 10

	jwtTokenStr, err := h.service.CreateJwt(r.Context(), tenYears, validatedCodeCheckReq)
	if err != nil {
		replyError(w, err, "create_token_failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    jwtTokenStr,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().AddDate(10, 0, 0), // 10 лет
		MaxAge:   86400 * 365 * 10,             // 10 лет
	})

	// удаляем код, иначе могут 1 успешным кодом создать очень много токенов
	if err := h.service.PurgeCode(r.Context(), validatedCodeCheckReq); err != nil {
		// просто логируем, не прерываем выполнение
		log.Printf("purgecode failed: %v", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "Token issued in cookie"}`))

}
