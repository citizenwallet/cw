package community

import (
	"math/big"
	"net/http"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/response"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/common/wei"
	"github.com/daobrussels/cw/pkg/cw"
	"github.com/daobrussels/smartcontracts/pkg/contracts/gateway"
	"github.com/daobrussels/smartcontracts/pkg/contracts/paymaster"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Handlers struct {
	supply     *supply.Supply
	ethservice *ethrequest.EthService
	responder  *response.Responder
	chain      *cw.ChainConfig
}

func NewHandlers(s *supply.Supply, e *ethrequest.EthService, r *response.Responder, chain *cw.ChainConfig) *Handlers {
	return &Handlers{
		s,
		e,
		r,
		chain,
	}
}

// GatewayDeploy returns the address of the deployed gateway contract
func (h *Handlers) Deploy(w http.ResponseWriter, r *http.Request) {
	// use the main wallet to deploy the gateway contract
	auth, err := bind.NewKeyedTransactorWithChainID(h.supply.PrivateKey, big.NewInt(int64(h.chain.ChainID)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get the next nonce for the main wallet
	nonce, err := h.ethservice.NextNonce(h.supply.Address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check the gas price
	gasPrice, err := h.ethservice.EstimateGasPrice(h.supply.Address, 0, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set default parameters
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = big.NewInt(int64(gasPrice))

	// deploy the gateway contract
	gaddr, _, _, err := gateway.DeployGateway(auth, h.ethservice.Client())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	print("gateway address:")
	print(gaddr.Hex())

	// deploy the paymaster
	// get the next nonce for the main wallet
	nonce, err = h.ethservice.NextNonce(h.supply.Address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auth.Nonce = big.NewInt(int64(nonce))

	paddr, _, p, err := paymaster.DeployPaymaster(auth, h.ethservice.Client(), gaddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	print("paymaster address:")
	print(paddr.Hex())

	// add initial deposit to the paymaster
	// get the next nonce for the main wallet
	nonce, err = h.ethservice.NextNonce(h.supply.Address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auth.Value = big.NewInt(int64(wei.WeiToEth(10)))

	_, err = p.Deposit(auth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	print("paymaster deposited:")
	dep, err := p.PaymasterCaller.GetDeposit(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	print(dep.String())

	// deploy the account factory

	// deploy the app factory

	err = h.responder.EncryptedBody(w, r.Context(), response.AddressResponse{Address: gaddr.Hex()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type SignedTx struct {
	TX string `json:"tx"`
}
