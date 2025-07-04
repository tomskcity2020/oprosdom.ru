package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"oprosdom.ru/monolith/internal/dz/internal/models"
	"oprosdom.ru/monolith/internal/dz/internal/repo"
	//"oprosdom.ru/monolith/internal/dz/internal/service"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME"))
}

func AddMember(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// парсим
	var member models.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// создаем uuid (отдельно делаем чтоб BasicValidate использовать и для POST и для PUT)
	if err := member.CreateUuid(); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
	}

	// проводим первичную проверку (структура заполнена полностью сейчас)
	if err := member.BasicValidate(); err != nil {
		replyError(w, err.Error(), http.StatusBadRequest)
		return
	}

	repository, err := repo.GetRepoSingleton()
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// сохраняем в файл
	if err := repository.SaveToFile(&member); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// сохраняем в слайс
	if err := repository.Save(&member); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// успех
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&member)

}

func replyError(w http.ResponseWriter, error string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": error,
	})
}

func AddKvartira(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// парсим
	var kvartira models.Kvartira
	if err := json.NewDecoder(r.Body).Decode(&kvartira); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// создаем uuid (отдельно делаем чтоб BasicValidate использовать и для POST и для PUT)
	if err := kvartira.CreateUuid(); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
	}

	// проводим первичную проверку (структура заполнена полностью сейчас)
	if err := kvartira.BasicValidate(); err != nil {
		replyError(w, err.Error(), http.StatusBadRequest)
		return
	}

	repository, err := repo.GetRepoSingleton()
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// сохраняем в файл
	if err := repository.SaveToFile(&kvartira); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// сохраняем в слайс
	if err := repository.Save(&kvartira); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// успех
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&kvartira)

}

func UpdateMember(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var member models.Member

	// извлекаем id, если нет id, то никакой ошибки или паники не будет, получим пустую строку ""
	varsMap := mux.Vars(r)
	idRaw := varsMap["id"]

	// добавляем пришедший id в структуру, пока не проверяем - проверка будет далее всех полей структуры сразу
	member.AddUuid(idRaw)

	// парсим данные в структуру
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// проводим первичную валидацию (когда структура заполнена включая id)
	if err := member.BasicValidate(); err != nil {
		replyError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// создаем репозиторий и вносим изменения в файл и соответствующий слайс
	repository, err := repo.GetRepoSingleton()
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repository.UpdateFile(&member); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repository.UpdateSlice(&member); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&member)

}

func UpdateKvartira(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var kvartira models.Kvartira

	// извлекаем id, если нет id, то никакой ошибки или паники не будет, получим пустую строку ""
	varsMap := mux.Vars(r)
	idRaw := varsMap["id"]

	// добавляем пришедший id в структуру, пока не проверяем - проверка будет далее всех полей структуры сразу
	kvartira.AddUuid(idRaw)

	// парсим данные в структуру
	if err := json.NewDecoder(r.Body).Decode(&kvartira); err != nil {
		replyError(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// проводим первичную валидацию (когда структура заполнена включая id)
	if err := kvartira.BasicValidate(); err != nil {
		replyError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// создаем репозиторий и вносим изменения в файл и соответствующий слайс
	repository, err := repo.GetRepoSingleton()
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repository.UpdateFile(&kvartira); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repository.UpdateSlice(&kvartira); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&kvartira)

}

func GetMembers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// обращаемся к репо, при обращении репозиторий подтянет данные из файла в слайс и затем выдаем данные слайса (который на текущий момент, он может в процессе пополнится)
	repository, err := repo.GetRepoSingleton()
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := repository.GetSliceMembers()

	json.NewEncoder(w).Encode(data)

}

func GetKvartiras(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// обращаемся к репо, при обращении репозиторий подтянет данные из файла в слайс и затем выдаем данные слайса (который на текущий момент, он может в процессе пополнится)
	repository, err := repo.GetRepoSingleton()
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := repository.GetSliceKvartiras()

	json.NewEncoder(w).Encode(data)

}

func GetMember(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	varsMap := mux.Vars(r)
	id := varsMap["id"]

	// чекаем что idRaw является корректным uuid, если не является - дальше не продолжаем, возвращаем ошибку
	_, err := uuid.Parse(id)
	if err != nil {
		replyError(w, "неправильный id жителя", http.StatusBadRequest)
		return
	}

	repository, err := repo.GetRepoSingleton()
	if err != nil {
		replyError(w, "ошибка репо", http.StatusInternalServerError)
		return
	}

	data, err := repository.GetMemberById(id)
	if err != nil {
		if err.Error() == "not_found" {
			replyError(w, "житель не найден", http.StatusNotFound)
		} else {
			replyError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(data)

}

func GetKvartira(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	varsMap := mux.Vars(r)
	id := varsMap["id"]

	// чекаем что idRaw является корректным uuid, если не является - дальше не продолжаем, возвращаем ошибку
	_, err := uuid.Parse(id)
	if err != nil {
		replyError(w, "неправильный id", http.StatusBadRequest)
		return
	}

	repository, err := repo.GetRepoSingleton()
	if err != nil {
		replyError(w, "ошибка репо", http.StatusInternalServerError)
		return
	}

	data, err := repository.GetKvartiraById(id)
	if err != nil {
		if err.Error() == "not_found" {
			replyError(w, "квартира не найдена", http.StatusNotFound)
		} else {
			replyError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(data)

}

func RemoveById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	varsMap := mux.Vars(r)
	mk := varsMap["mk"]
	id := varsMap["id"]

	var filename string

	switch mk {
	case "member":
		filename = "members"
	case "kvartira":
		filename = "kvartiras"
	default:
		replyError(w, "неправильный запрос", http.StatusBadRequest)
		return
	}

	// чекаем что idRaw является корректным uuid, если не является - дальше не продолжаем, возвращаем ошибку
	_, err := uuid.Parse(id)
	if err != nil {
		replyError(w, "неправильный id", http.StatusBadRequest)
		return
	}

	repository, err := repo.GetRepoSingleton()
	if err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repository.RemoveFromFile(filename, id); err != nil {
		replyError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch filename {
	case "members":
		if err := repository.RemoveMemberSlice(id); err != nil {
			replyError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "kvartiras":
		if err := repository.RemoveKvartiraSlice(id); err != nil {
			replyError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		replyError(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
}
