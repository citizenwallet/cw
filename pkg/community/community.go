package community

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/smartcontracts/pkg/contracts/accfactory"
	"github.com/daobrussels/smartcontracts/pkg/contracts/gateway"
	"github.com/daobrussels/smartcontracts/pkg/contracts/grfactory"
	"github.com/daobrussels/smartcontracts/pkg/contracts/paymaster"
	"github.com/daobrussels/smartcontracts/pkg/contracts/profactory"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func DeployGateway(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chainID *big.Int) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address.Hex())
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

func DeployPaymaster(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chainID *big.Int, entryPoint common.Address) (*common.Address, *paymaster.Paymaster, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address.Hex())
	if err != nil {
		return nil, nil, err
	}

	// set default parameters
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// deploy the paymaster contract
	addr, _, p, err := paymaster.DeployPaymaster(auth, es.Client(), entryPoint)
	if err != nil {
		return nil, nil, err
	}

	return &addr, p, nil
}

func FundPaymaster(p *paymaster.Paymaster, es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chainID *big.Int, amount *big.Int) error {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address.Hex())
	if err != nil {
		return err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = amount

	_, err = p.Deposit(auth)
	if err != nil {
		return err
	}

	return nil
}

func DeployAccountFactory(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chainID *big.Int, entryPoint common.Address) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address.Hex())
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

func DeployGratitudeFactory(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chainID *big.Int, entryPoint common.Address) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address.Hex())
	if err != nil {
		return nil, err
	}

	// set default parameters
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// deploy the gratitude factory contract
	addr, _, _, err := grfactory.DeployGrfactory(auth, es.Client(), entryPoint)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func CreateGratitudeApp(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chainID *big.Int, owner, faddr common.Address) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address.Hex())
	if err != nil {
		return nil, err
	}

	// set default parameters
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// instantiate gratitude factory contract
	factory, err := grfactory.NewGrfactory(faddr, es.Client())
	if err != nil {
		return nil, err
	}

	// create the gratitude app
	_, err = factory.CreateGratitudeToken(auth, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	addr, err := factory.GetGratitudeTokenAddress(&bind.CallOpts{}, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func CreateAccount(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chainID *big.Int, owner, faddr common.Address) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address.Hex())
	if err != nil {
		return nil, err
	}

	// set default parameters
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// instantiate account factory contract
	factory, err := accfactory.NewAccfactory(faddr, es.Client())
	if err != nil {
		return nil, err
	}

	_, err = factory.CreateAccount(auth, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	addr, err := factory.GetAddress(&bind.CallOpts{}, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func DeployProfileFactory(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chainID *big.Int, entryPoint common.Address) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address.Hex())
	if err != nil {
		return nil, err
	}

	// set default parameters
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// deploy profile factory contract
	addr, _, _, err := profactory.DeployProfactory(auth, es.Client(), entryPoint)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func CreateProfile(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chainID *big.Int, owner, faddr common.Address) (*common.Address, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := es.NextNonce(address.Hex())
	if err != nil {
		return nil, err
	}

	// set default parameters
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)

	// instantiate profile factory contract
	factory, err := profactory.NewProfactory(faddr, es.Client())
	if err != nil {
		return nil, err
	}

	_, err = factory.CreateProfile(auth, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	addr, err := factory.GetProfileAddress(&bind.CallOpts{}, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	return &addr, nil
}
