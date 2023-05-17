package ethrequest

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	ETHEstimateGas        = "eth_estimateGas"
	ETHSendRawTransaction = "eth_sendRawTransaction"
)

type EthService struct {
	rpc    *rpc.Client
	client *ethclient.Client
	ctx    context.Context
}

func (e *EthService) Client() *ethclient.Client {
	return e.client
}

func NewEthService(endpoint string) (*EthService, error) {
	rpc, err := rpc.DialHTTP(endpoint)
	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(rpc)

	return &EthService{rpc, client, context.Background()}, nil
}

func (e *EthService) Close() {
	e.ctx.Done()
	e.client.Close()
}

func (e *EthService) EstimateGas(from, to string, value uint64) (uint64, error) {
	t := common.HexToAddress(to)

	msg := ethereum.CallMsg{
		From:  common.HexToAddress(from),
		To:    &t,
		Value: big.NewInt(int64(value)),
		Gas:   0,
	}

	return e.client.EstimateGas(e.ctx, msg)
}

func (e *EthService) EstimateGasPrice(from string, value uint64, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From:  common.HexToAddress(from),
		Value: big.NewInt(int64(value)),
		Data:  data,
		Gas:   0,
	}

	return e.client.EstimateGas(e.ctx, msg)
}

func (e *EthService) EstimateContractGasPrice(data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		Data: data,
		Gas:  0,
	}

	return e.client.EstimateGas(e.ctx, msg)
}

func (e *EthService) SendRawTransaction(tx string) ([]byte, error) {

	err := e.rpc.Call(nil, ETHSendRawTransaction, tx)

	return nil, err
}

func (e *EthService) NextNonce(address string) (uint64, error) {
	return e.client.PendingNonceAt(e.ctx, common.HexToAddress(address))
}
