package transaction

import "net/http"

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Send(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
