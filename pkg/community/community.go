package community

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/cw"
	"github.com/daobrussels/smartcontracts/pkg/contracts/accfactory"
	"github.com/daobrussels/smartcontracts/pkg/contracts/account"
	"github.com/daobrussels/smartcontracts/pkg/contracts/gateway"
	"github.com/daobrussels/smartcontracts/pkg/contracts/grfactory"
	"github.com/daobrussels/smartcontracts/pkg/contracts/paymaster"
	"github.com/daobrussels/smartcontracts/pkg/contracts/profactory"
	"github.com/daobrussels/smartcontracts/pkg/contracts/profile"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type CommunityAddress struct {
	Gateway          common.Address `json:"gateway"`
	Paymaster        common.Address `json:"paymaster"`
	AccountFactory   common.Address `json:"accountFactory"`
	GratitudeFactory common.Address `json:"gratitudeFactory"`
	ProfileFactory   common.Address `json:"profileFactory"`
	Chain            cw.ChainConfig `json:"chain"`
}

type Community struct {
	es      *ethrequest.EthService
	key     *ecdsa.PrivateKey
	address common.Address
	chain   cw.ChainConfig

	EntryPoint common.Address
	Gateway    *gateway.Gateway

	paddr     common.Address
	Paymaster *paymaster.Paymaster

	afaddr         common.Address
	AccountFactory *accfactory.Accfactory

	grfaddr          common.Address
	GratitudeFactory *grfactory.Grfactory

	prfaddr        common.Address
	ProfileFactory *profactory.Profactory
}

func (c *Community) ExportAddress() CommunityAddress {
	return CommunityAddress{
		Gateway:          c.EntryPoint,
		Paymaster:        c.paddr,
		AccountFactory:   c.afaddr,
		GratitudeFactory: c.grfaddr,
		ProfileFactory:   c.prfaddr,
		Chain:            c.chain,
	}
}

// New instantiates a community struct using the provided addresses for the contracts
func New(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chain cw.ChainConfig, addr CommunityAddress) (*Community, error) {
	// instantiate gateway contract
	g, err := gateway.NewGateway(addr.Gateway, es.Client())
	if err != nil {
		return nil, err
	}

	// instantiate paymaster contract
	p, err := paymaster.NewPaymaster(addr.Paymaster, es.Client())
	if err != nil {
		return nil, err
	}

	// instantiate account factory contract
	acc, err := accfactory.NewAccfactory(addr.AccountFactory, es.Client())
	if err != nil {
		return nil, err
	}

	// instantiate gratitude factory contract
	gr, err := grfactory.NewGrfactory(addr.GratitudeFactory, es.Client())
	if err != nil {
		return nil, err
	}

	// instantiate profile factory contract
	pro, err := profactory.NewProfactory(addr.ProfileFactory, es.Client())
	if err != nil {
		return nil, err
	}

	return &Community{
		es:               es,
		key:              key,
		address:          address,
		chain:            chain,
		EntryPoint:       addr.Gateway,
		Gateway:          g,
		paddr:            addr.Paymaster,
		Paymaster:        p,
		afaddr:           addr.AccountFactory,
		AccountFactory:   acc,
		grfaddr:          addr.GratitudeFactory,
		GratitudeFactory: gr,
		prfaddr:          addr.ProfileFactory,
		ProfileFactory:   pro,
	}, nil
}

// Deploy instantiates a community struct and deploys the contracts
func Deploy(es *ethrequest.EthService, key *ecdsa.PrivateKey, address common.Address, chain cw.ChainConfig) (*Community, error) {
	c := &Community{
		es:      es,
		key:     key,
		address: address,
		chain:   chain,
	}

	// instantiate gateway contract
	err := c.DeployGateway()
	if err != nil {
		return nil, err
	}

	// deploy paymaster contract
	err = c.DeployPaymaster()
	if err != nil {
		return nil, err
	}

	// deploy account factory contract
	err = c.DeployAccountFactory()
	if err != nil {
		return nil, err
	}

	// deploy gratitude factory contract
	err = c.DeployGratitudeFactory()
	if err != nil {
		return nil, err
	}

	// deploy profile factory contract
	err = c.DeployProfileFactory()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewTransactor returns a new transactor for the community
func (c *Community) NewTransactor() (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(c.key, big.NewInt(int64(c.chain.ChainID)))
}

// NextNonce returns the next nonce for the community
func (c *Community) NextNonce() (uint64, error) {
	return c.es.NextNonce(c.address.Hex())
}

// DeployGateway deploys the gateway contract
func (c *Community) DeployGateway() error {
	auth, err := c.NewTransactor()
	if err != nil {
		return err
	}

	// get the next nonce for the main wallet
	nonce, err := c.NextNonce()
	if err != nil {
		return err
	}

	// set default parameters
	setDefaultParameters(auth, nonce)

	// deploy the gateway contract
	addr, _, g, err := gateway.DeployGateway(auth, c.es.Client())
	if err != nil {
		return err
	}

	c.EntryPoint = addr
	c.Gateway = g

	return nil
}

// DeployPaymaster deploys the paymaster contract
func (c *Community) DeployPaymaster() error {
	auth, err := c.NewTransactor()
	if err != nil {
		return err
	}

	// get the next nonce for the main wallet
	nonce, err := c.NextNonce()
	if err != nil {
		return err
	}

	// set default parameters
	setDefaultParameters(auth, nonce)

	// deploy the paymaster contract
	addr, _, p, err := paymaster.DeployPaymaster(auth, c.es.Client(), c.EntryPoint)
	if err != nil {
		return err
	}

	c.paddr = addr
	c.Paymaster = p

	return nil
}

// FundPaymaster funds the paymaster contract
func (c *Community) FundPaymaster(amount *big.Int) error {
	auth, err := c.NewTransactor()
	if err != nil {
		return err
	}

	// get the next nonce for the main wallet
	nonce, err := c.NextNonce()
	if err != nil {
		return err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = amount

	_, err = c.Paymaster.Deposit(auth)
	if err != nil {
		return err
	}

	return nil
}

// DeployAccountFactory deploys the account factory contract
func (c *Community) DeployAccountFactory() error {
	auth, err := c.NewTransactor()
	if err != nil {
		return err
	}

	// get the next nonce for the main wallet
	nonce, err := c.NextNonce()
	if err != nil {
		return err
	}

	// set default parameters
	setDefaultParameters(auth, nonce)

	// deploy the account factory contract
	addr, _, acc, err := accfactory.DeployAccfactory(auth, c.es.Client(), c.EntryPoint)
	if err != nil {
		return err
	}

	c.afaddr = addr
	c.AccountFactory = acc

	return nil
}

// DeployGratitudeFactory deploys the gratitude factory contract
func (c *Community) DeployGratitudeFactory() error {
	auth, err := c.NewTransactor()
	if err != nil {
		return err
	}

	// get the next nonce for the main wallet
	nonce, err := c.NextNonce()
	if err != nil {
		return err
	}

	// set default parameters
	setDefaultParameters(auth, nonce)

	// deploy the gratitude factory contract
	addr, _, gr, err := grfactory.DeployGrfactory(auth, c.es.Client(), c.EntryPoint)
	if err != nil {
		return err
	}

	c.grfaddr = addr
	c.GratitudeFactory = gr

	return nil
}

// CreateGratitudeApp creates a gratitude app for the provided owner
func (c *Community) CreateGratitudeApp(owner common.Address) (*common.Address, error) {
	auth, err := c.NewTransactor()
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := c.NextNonce()
	if err != nil {
		return nil, err
	}

	// set default parameters
	setDefaultParameters(auth, nonce)

	// create the gratitude app
	_, err = c.GratitudeFactory.CreateGratitudeToken(auth, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	addr, err := c.GratitudeFactory.GetGratitudeTokenAddress(&bind.CallOpts{}, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

// CreateAccount creates an account for the provided owner
func (c *Community) CreateAccount(owner common.Address) (*common.Address, error) {
	auth, err := c.NewTransactor()
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := c.NextNonce()
	if err != nil {
		return nil, err
	}

	// set default parameters
	setDefaultParameters(auth, nonce)

	_, err = c.AccountFactory.CreateAccount(auth, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	addr, err := c.AccountFactory.GetAddress(&bind.CallOpts{}, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

// DeployProfileFactory deploys the profile factory contract
func (c *Community) DeployProfileFactory() error {
	auth, err := c.NewTransactor()
	if err != nil {
		return err
	}

	// get the next nonce for the main wallet
	nonce, err := c.NextNonce()
	if err != nil {
		return err
	}

	// set default parameters
	setDefaultParameters(auth, nonce)

	// deploy profile factory contract
	addr, _, pr, err := profactory.DeployProfactory(auth, c.es.Client(), c.EntryPoint)
	if err != nil {
		return err
	}

	c.prfaddr = addr
	c.ProfileFactory = pr

	return nil
}

// CreateProfile creates a profile for the provided owner
func (c *Community) CreateProfile(owner common.Address) (*common.Address, error) {
	auth, err := c.NewTransactor()
	if err != nil {
		return nil, err
	}

	// get the next nonce for the main wallet
	nonce, err := c.NextNonce()
	if err != nil {
		return nil, err
	}

	// set default parameters
	setDefaultParameters(auth, nonce)

	_, err = c.ProfileFactory.CreateProfile(auth, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	addr, err := c.ProfileFactory.GetProfileAddress(&bind.CallOpts{}, owner, big.NewInt(int64(nonce)))
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

// GetProfile returns the profile for the provided owner
func (c *Community) GetProfile(owner common.Address) (*profile.Profile, error) {
	p, err := profile.NewProfile(owner, c.es.Client())
	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetAccount returns the account for the provided owner
func (c *Community) GetAccount(owner common.Address) (*account.Account, error) {
	a, err := account.NewAccount(owner, c.es.Client())
	if err != nil {
		return nil, err
	}

	return a, nil
}

// setDefaultParameters sets the nonce, value and gas limit for a default contract transaction
func setDefaultParameters(auth *bind.TransactOpts, nonce uint64) {
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000 - 1)
}
