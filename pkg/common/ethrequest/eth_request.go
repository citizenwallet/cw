package ethrequest

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	ETHEstimateGas        = "eth_estimateGas"
	ETHSendRawTransaction = "eth_sendRawTransaction"

	ETHSendUserOperation   = "eth_sendUserOperation"
	PMSponsorUserOperation = "pm_sponsorUserOperation"
)

type UserOperation struct {
	Sender               string `json:"sender"`
	Nonce                string `json:"nonce"`
	InitCode             string `json:"initCode"`
	CallData             string `json:"callData"`
	CallGasLimit         string `json:"callGasLimit"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	PreVerificationGas   string `json:"preVerificationGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	PaymasterAndData     string `json:"paymasterAndData"`
	Signature            string `json:"signature"`
}

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

func (e *EthService) EstimateFullGas(from common.Address, tx *types.Transaction) (uint64, error) {

	msg := ethereum.CallMsg{
		From:       from,
		To:         tx.To(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		GasFeeCap:  tx.GasFeeCap(),
		GasTipCap:  tx.GasTipCap(),
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	}

	return e.client.EstimateGas(e.ctx, msg)
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

func (e *EthService) GetCode(address common.Address) ([]byte, error) {
	return e.client.CodeAt(e.ctx, address, nil)
}

type SponsorOp struct {
	NextNonce            string `json:"nextNonce"`
	PaymasterAndData     string `json:"paymasterAndData"`
	PreVerificationGas   string `json:"preVerificationGas"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	CallGasLimit         string `json:"callGasLimit"`
}

func (e *EthService) SponsorUserOp(nonce string, op []byte, eaddr, ptype string) (*SponsorOp, error) {
	uop := &UserOperation{}
	err := json.Unmarshal(op, uop)
	if err != nil {
		return nil, err
	}

	uop.Nonce = nonce

	println("nonce: ", uop.Nonce)

	sop := &SponsorOp{}

	err = e.rpc.Call(sop, PMSponsorUserOperation, uop, eaddr, map[string]string{"type": ptype})
	if err != nil {
		println(err.Error())
		return nil, err
	}

	sop.NextNonce = nonce
	sop.PreVerificationGas = makeValidEvenHex(sop.PreVerificationGas)
	sop.VerificationGasLimit = makeValidEvenHex(sop.VerificationGasLimit)
	sop.CallGasLimit = makeValidEvenHex(sop.CallGasLimit)

	return sop, nil
}

func makeValidEvenHex(h string) string {
	h = strip0x(h)
	h = evenHex(h)
	return "0x" + h
}

func strip0x(h string) string {
	if len(h) > 2 && h[:2] == "0x" {
		return h[2:]
	}

	return h
}

func evenHex(h string) string {
	if len(h)%2 == 0 {
		return h
	}

	return "0" + h
}

func (e *EthService) SendUserOp(op []byte, eaddr string) error {
	uop := &UserOperation{}
	err := json.Unmarshal(op, uop)
	if err != nil {
		return err
	}

	err = e.rpc.Call(nil, ETHSendUserOperation, uop, eaddr)
	if err != nil {
		println(err.Error())
		return err
	}

	return nil
}
