package ethrequest

import (
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	ETHEstimateGas        = "eth_estimateGas"
	ETHSendRawTransaction = "eth_sendRawTransaction"
)

type EthService struct {
	client  *rpc.Client
	rClient *RawService
}

func NewEthService(endpoint string) (*EthService, error) {
	client, err := rpc.DialHTTP(endpoint)
	if err != nil {
		return nil, err
	}

	rclient := NewRawService(endpoint)

	return &EthService{client, rclient}, nil
}

func (e *EthService) Close() {
	e.client.Close()
}

func (e *EthService) EstimateGas(txobj any) (string, error) {

	gas := "0x0"

	err := e.client.Call(&gas, ETHEstimateGas, txobj, "latest")

	return gas, err
}

func (e *EthService) SendRawTransaction(tx string) ([]byte, error) {

	err := e.client.Call(nil, ETHSendRawTransaction, tx)

	return nil, err
}
