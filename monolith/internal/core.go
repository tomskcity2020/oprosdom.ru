package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestJsonRpc struct {
	JsonRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Id      int                    `json:"id"`
}

type ResponseJsonRpc struct {
	JsonRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"` // результата не будет если есть ошибка
	Error   string      `json:"error,omitempty"`  // ошибки не будет если есть результат
	Id      int         `json:"id,omitempty"`     // если id нет, то это значит чтоне удалось распарсить инициальный запрос, то есть id мы не знаем
}

type JsonRpcHandler struct{}

func NewJsonRpcHandler() *JsonRpcHandler {
	return &JsonRpcHandler{}
}

func (obj *JsonRpcHandler) RequestResponse(w http.ResponseWriter, r *http.Request) {

	// используем decode вместо unmarshal потому, что decode постепенно считывает данные, а unmarshal загружает сначала весь байтовый срез в память
	decoder := json.NewDecoder(r.Body)

	var rpcRequest RequestJsonRpc
	err := decoder.Decode(&rpcRequest)
	if err != nil {
		http.Error(w, "Incorrect JSON-RPC structure or decoding error", http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v", rpcRequest)

	// по методу создаем обработчик (вернется указатель на структуру соответствующую методу)
	handlerStruct, exists := Handlers[rpcRequest.Method]
	if !exists {
		http.Error(w, "Unknown method", http.StatusNotFound)
		return
	}

	response, err := handlerStruct.Handle(rpcRequest.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rcpResponse := ResponseJsonRpc{
		JsonRPC: "2.0",
		Result:  response,
		Id:      rpcRequest.Id,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rcpResponse); err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding failed: %s", err), http.StatusInternalServerError)
		return
	}

}
