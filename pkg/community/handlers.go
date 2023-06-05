package community

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strings"

	"github.com/daobrussels/cw/pkg/common/response"
	"github.com/daobrussels/cw/pkg/cw"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
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

	println(addr)

	acc, err := h.c.CreateAccount(common.HexToAddress(addr))
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.responder.Body(w, response.AddressResponse{Address: acc.Hex()})
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type balanceResp struct {
	Balance *big.Int `json:"balance"`
}

// GetAccountDERC20Balance returns the balance of a given account
func (h *Handlers) GetAccountDERC20Balance(w http.ResponseWriter, r *http.Request) {
	accAddr := chi.URLParam(r, "account_id")
	if accAddr == "" || accAddr == "0x" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bal, err := h.c.GetDERC20Balance(common.HexToAddress(accAddr))
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.responder.Body(w, balanceResp{Balance: bal})
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// CreateProfile creates a profile in the community for a given account and returns the address
func (h *Handlers) CreateProfile(w http.ResponseWriter, r *http.Request) {
	accAddr := chi.URLParam(r, "account_id")
	if accAddr == "" || accAddr == "0x" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addr, ok := cw.GetAddressFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	acc, err := h.c.GetAccount(common.HexToAddress(accAddr))
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	owner, err := acc.Owner(&bind.CallOpts{})
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	println(addr)

	if strings.ToLower(owner.Hex()) != strings.ToLower(addr) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	println(accAddr)

	pr, err := h.c.CreateProfile(common.HexToAddress(accAddr))
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.responder.Body(w, response.AddressResponse{Address: pr.Hex()})
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type opRequest struct {
	Sender string `json:"sender"`
	Op     string `json:"op"`
}

// SubmitOp submits an operation to the gateway for processing
func (h *Handlers) SubmitOp(w http.ResponseWriter, r *http.Request) {
	addr, ok := cw.GetAddressFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	println("addr " + addr)

	req := &opRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		println("decoding ")
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	println("data " + req.Op)

	err = h.c.SubmitOp([]byte(req.Op))
	if err != nil {
		println("op")
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type voucherRequest struct {
	Amount      int64  `json:"amount"`
	MinterName  string `json:"minter_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateVoucher creates a voucher in the community
func (h *Handlers) CreateVoucher(w http.ResponseWriter, r *http.Request) {
	addr, ok := cw.GetAddressFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := &voucherRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	meta, err := h.c.RegensMintCustomVoucher(common.HexToAddress(addr), big.NewInt(1), req.MinterName, req.Name, req.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.responder.Body(w, meta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetVouchers returns vouchers from the community
func (h *Handlers) GetVouchers(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	owners := []common.Address{}
	froms := []common.Address{}
	tos := []common.Address{}

	for k, v := range values {
		switch k {
		case "owner":
			owners = append(owners, common.HexToAddress(v[0]))
		case "from":
			froms = append(froms, common.HexToAddress(v[0]))
		case "to":
			tos = append(tos, common.HexToAddress(v[0]))
		}
	}

	vouchers, err := h.c.RegensFilterTransferSingle(owners, froms, tos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.responder.BodyMultiple(w, vouchers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Sponsor returns the paymaster address for the community along with gas data
func (h *Handlers) SponsorOp(w http.ResponseWriter, r *http.Request) {
	opr := &opRequest{}

	err := json.NewDecoder(r.Body).Decode(opr)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	println("op")
	println(opr.Op)

	sop, err := h.c.GetPaymasterData(common.HexToAddress(opr.Sender), []byte(opr.Op))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.responder.Body(w, sop)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
