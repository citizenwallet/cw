package transaction

import (
	"encoding/json"
	"net/http"
)

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

type SignedTx struct {
	TX string `json:"tx"`
}

func (h *Handlers) Send(w http.ResponseWriter, r *http.Request) {
	var req SignedTx

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	println(req.TX)
}
