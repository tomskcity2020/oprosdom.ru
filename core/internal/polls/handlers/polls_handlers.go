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
