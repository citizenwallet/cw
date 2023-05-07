package hello

import (
	"net/http"

	"github.com/daobrussels/cw/pkg/common/response"
)

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

// Hello returns the supply wallet public key and address
func (h *Handlers) Hello(w http.ResponseWriter, r *http.Request) {
	response.EncryptedBody(w, r.Context(), []byte("{}"))
}
