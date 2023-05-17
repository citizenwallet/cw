package transaction

import (
	"fmt"
	"math/big"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/cw"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Service struct {
	chain      *cw.ChainConfig
	supply     *supply.Supply
	ethservice *ethrequest.EthService
}

func New(chain *cw.ChainConfig, s *supply.Supply, ethservice *ethrequest.EthService) *Service {
	return &Service{
		chain,
		s,
		ethservice,
	}
}

func (s *Service) Send(to string, amount int64) error {
	address := common.HexToAddress(to)

	gas, err := s.ethservice.EstimateGas(s.supply.Address, to, uint64(amount))
	if err != nil {
		return err
	}

	txdata := types.LegacyTx{
		Gas:      gas,
		GasPrice: big.NewInt(1000000000),
		To:       &address,
		Value:    big.NewInt(amount),
		Data:     nil,
	}

	sign := types.NewEIP155Signer(big.NewInt(int64(s.chain.ChainID)))

	tx, err := types.SignNewTx(s.supply.PrivateKey, sign, &txdata)
	if err != nil {
		return err
	}

	btx, err := tx.MarshalBinary()
	if err != nil {
		return err
	}

	paddedTx := fmt.Sprintf("0x%s", common.Bytes2Hex(btx))

	_, err = s.ethservice.SendRawTransaction(paddedTx)

	return err
}

func (s *Service) Forward(tx string) error {
	ethservice, err := ethrequest.NewEthService(s.chain.RPC[0])
	if err != nil {
		return err
	}
	defer ethservice.Close()

	paddedTx := fmt.Sprintf("0x%s", tx)

	_, err = ethservice.SendRawTransaction(paddedTx)

	return err
}
