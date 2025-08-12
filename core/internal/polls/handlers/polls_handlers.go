package polls_handlers

import (
	"encoding/json"
	"log"
	"net/http"

	polls_models "oprosdom.ru/core/internal/polls/models"
	polls_service "oprosdom.ru/core/internal/polls/service"
)

type Handler struct {
	service polls_service.PollsService
}

func NewHandler(service polls_service.PollsService) *Handler {
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

// GetPolls godoc
// @Summary      Получить список опросов
// @Description  Возвращает список доступных опросов с их ID и заголовками. Требует JWT в cookie 'auth'.
// @Tags         polls
// @Produce      json
// @Success      200  {array}   polls_models.Poll   "Список опросов"
// @Failure      401  {string}  string              "unauthorized — JWT не предоставлен или неверен"
// @Failure      500  {string}  string              "something wrong with polls service"
// @Security     CookieAuth
// @Router       /polls [get]
func (h *Handler) GetPolls(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// service (там репо и бизнес)
	data, err := h.service.GetPolls(r.Context())
	if err != nil {
		replyError(w, err, "something wrong with polls service", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)

}

// Vote godoc
// @Summary      Проголосовать в опросе
// @Description  Отправляет голос пользователя за или против по конкретному опросу. Требует JWT в cookie 'auth'.
// @Tags         polls
// @Accept       json
// @Produce      json
// @Param        vote  body      polls_models.UnsafeVoteReq  true  "JSON с ID опроса и голосом"
// @Success      200   {string}  string  "success"
// @Failure      400   {string}  string  "incorrect_vote_request / vote_validation_failed"
// @Failure      401   {string}  string  "unauthorized — JWT не предоставлен или неверен"
// @Failure      500   {string}  string  "vote_wrong_service"
// @Security     CookieAuth
// @Router       /polls/vote [post]
func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var unsafeVoteReq polls_models.UnsafeVoteReq
	if err := json.NewDecoder(r.Body).Decode(&unsafeVoteReq); err != nil {
		replyError(w, err, "incorrect_vote_request", http.StatusBadRequest)
		return
	}

	validVoteReq, err := unsafeVoteReq.Validate()
	if err != nil {
		replyError(w, err, "vote_validation_failed", http.StatusBadRequest)
		return
	}

	// получаем jti из контекста!!!
	jti, ok := r.Context().Value("jti").(string)
	if !ok {
		replyError(w, err, "jti_wrong", http.StatusInternalServerError)
		return
	}

	validVoteReq.Jti = jti

	if err := h.service.Vote(r.Context(), validVoteReq); err != nil {
		replyError(w, err, "vote_wrong_service", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("success")

}

// PollStats godoc
// @Summary      Получить статистику по опросам
// @Description  Возвращает статистику голосования по каждому опросу. Требует JWT в cookie 'auth'.
// @Tags         polls
// @Produce      json
// @Success      200  {array}   polls_models.PollStats  "Массив статистики по опросам"
// @Failure      401  {string}  string                  "unauthorized — JWT не предоставлен или неверен"
// @Failure      500  {string}  string                  "something wrong with polls service"
// @Security     CookieAuth
// @Router       /polls/stats [get]
func (h *Handler) PollStats(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// service (там репо и бизнес)
	data, err := h.service.PollStats(r.Context())
	if err != nil {
		replyError(w, err, "something wrong with polls service", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)

}
