package transaction

import (
	"encoding/json"
	"net/http"

	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/common/transaction"
	"github.com/daobrussels/cw/pkg/cw"
)

type Handlers struct {
	tr *transaction.Service
}

func NewHandlers(chain *cw.ChainConfig,
	supply *supply.Supply) *Handlers {
	return &Handlers{
		tr: transaction.New(chain, supply),
	}
}

type SignedTx struct {
	TX string `json:"tx"`
}

func (h *Handlers) Send(w http.ResponseWriter, r *http.Request) {
	var req SignedTx

	println("send")

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	println("sending...")

	err = h.tr.Forward(req.TX)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	println("sent!!!")
}
