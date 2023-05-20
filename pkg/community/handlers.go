package community

import (
	"net/http"

	"github.com/daobrussels/cw/pkg/common/response"
)

type Handlers struct {
	responder *response.Responder
	c         *Community
}

func NewHandlers(r *response.Responder, c *Community) *Handlers {
	return &Handlers{
		r,
		c,
	}
}

// Config returns the community config of addresses and chain info
func (h *Handlers) Config(w http.ResponseWriter, r *http.Request) {
	addr := h.c.ExportAddress()

	err := h.responder.EncryptedBody(w, r.Context(), addr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
