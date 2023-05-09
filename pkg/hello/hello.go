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
	chain     cw.ChainConfig
	responder *response.Responder
}

func NewHandlers(chain cw.ChainConfig, r *response.Responder) *Handlers {
	return &Handlers{
		chain:     chain,
		responder: r,
	}
}

// Hello returns returns the local chain configuration and signs the response.
// Allows for clients to verify the response and respond using the public key of the sender
func (h *Handlers) Hello(w http.ResponseWriter, r *http.Request) {
	err := h.responder.EncryptedBody(w, r.Context(), h.chain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
