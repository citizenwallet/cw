package push

import "net/http"

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Associate(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *Handlers) Dissociate(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
