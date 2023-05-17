package community

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/smartcontracts/pkg/contracts/accfactory"
	"github.com/daobrussels/smartcontracts/pkg/contracts/appfactory"
	"github.com/daobrussels/smartcontracts/pkg/contracts/gateway"
	"github.com/daobrussels/smartcontracts/pkg/contracts/paymaster"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func DeployGateway(es *ethrequest.EthService, key *ecdsa.PrivateKey, address string, chainID *big.Int) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address)
	if err != nil {
		return nil, err
	}

	// set default parameters
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// deploy the gateway contract
	addr, _, _, err := gateway.DeployGateway(auth, es.Client())
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func DeployPaymaster(es *ethrequest.EthService, key *ecdsa.PrivateKey, address string, chainID *big.Int, entryPoint common.Address) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address)
	if err != nil {
		return nil, err
	}

	// set default parameters
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// deploy the paymaster contract
	addr, _, _, err := paymaster.DeployPaymaster(auth, es.Client(), entryPoint)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func DeployAccountFactory(es *ethrequest.EthService, key *ecdsa.PrivateKey, address string, chainID *big.Int, entryPoint common.Address) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address)
	if err != nil {
		return nil, err
	}

	// set default parametersß
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// deploy the account factory contract
	addr, _, _, err := accfactory.DeployAccfactory(auth, es.Client(), entryPoint)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func DeployAppFactory(es *ethrequest.EthService, key *ecdsa.PrivateKey, address string, chainID *big.Int, entryPoint common.Address) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address)
	if err != nil {
		return nil, err
	}

	// set default parametersß
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// deploy the app factory contract
	addr, _, _, err := appfactory.DeployAppfactory(auth, es.Client(), entryPoint)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}
