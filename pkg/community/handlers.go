package community

import (
	"encoding/json"
	"net/http"

	"github.com/daobrussels/cw/pkg/common/response"
	"github.com/daobrussels/cw/pkg/cw"
	"github.com/ethereum/go-ethereum/common"
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

// CreateAccount creates an account in the community and returns the address
func (h *Handlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	addr, ok := cw.GetAddressFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	acc, err := h.c.CreateAccount(common.HexToAddress(addr))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.responder.EncryptedBody(w, r.Context(), response.AddressResponse{Address: acc.Hex()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type SubmitOpRequest struct {
	Data []byte `json:"data"`
}

// SubmitOp submits an operation to the gateway for processing
func (h *Handlers) SubmitOp(w http.ResponseWriter, r *http.Request) {
	addr, ok := cw.GetAddressFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req *SubmitOpRequest

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = h.c.SubmitOp(common.HexToAddress(addr), req.Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
