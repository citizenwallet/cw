package hello

import (
	"net/http"

	"github.com/daobrussels/cw/pkg/common/response"
	"github.com/daobrussels/cw/pkg/cw"
)

type HelloResponse struct {
	Message string `json:"message"`
}

type Handlers struct {
	responder *response.Responder
}

func NewHandlers(r *response.Responder) *Handlers {
	return &Handlers{
		responder: r,
	}
}

// Hello returns returns the local chain configuration and signs the response.
// Allows for clients to verify the response and respond using the public key of the sender
func (h *Handlers) Hello(w http.ResponseWriter, r *http.Request) {
	chain, err := cw.GetChain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.responder.EncryptedBody(w, r.Context(), chain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
