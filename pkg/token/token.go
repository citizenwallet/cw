package token

import "net/http"

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Mint(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *Handlers) Burn(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
