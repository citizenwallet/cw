package hello

import (
	"net/http"

	"github.com/daobrussels/cw/pkg/common/response"
)

type Handlers struct {
	responder *response.Responder
}

func NewHandlers(r *response.Responder) *Handlers {
	return &Handlers{
		responder: r,
	}
}

// Hello returns the supply wallet public key and address
func (h *Handlers) Hello(w http.ResponseWriter, r *http.Request) {
	err := h.responder.EncryptedBody(w, r.Context(), []byte("{}"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
