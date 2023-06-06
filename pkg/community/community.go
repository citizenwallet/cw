package community

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/voucher"
	"github.com/daobrussels/cw/pkg/common/wei"
	"github.com/daobrussels/cw/pkg/cw"
	"github.com/daobrussels/smartcontracts/pkg/contracts/account"
	"github.com/daobrussels/smartcontracts/pkg/contracts/derc20"
	"github.com/daobrussels/smartcontracts/pkg/contracts/gateway"
	"github.com/daobrussels/smartcontracts/pkg/contracts/grfactory"
	"github.com/daobrussels/smartcontracts/pkg/contracts/paymaster"
	"github.com/daobrussels/smartcontracts/pkg/contracts/profactory"
	"github.com/daobrussels/smartcontracts/pkg/contracts/profile"
	"github.com/daobrussels/smartcontracts/pkg/contracts/regensToken"
	"github.com/daobrussels/smartcontracts/pkg/contracts/simpleaccountfactory"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type CommunityAddress struct {
	Gateway          common.Address `json:"gateway"`
	Paymaster        common.Address `json:"paymaster"`
	AccountFactory   common.Address `json:"accountFactory"`
	GratitudeFactory common.Address `json:"gratitudeFactory"`
	ProfileFactory   common.Address `json:"profileFactory"`
	RegensToken      common.Address `json:"regensToken"`
	DERC20           common.Address `json:"derc20"`
	Chain            cw.ChainConfig `json:"chain"`
}

type PaymasterType string

const (
	PaymasterTypePayAsYouGo PaymasterType = "payg"
	PaymasterTypeToken      PaymasterType = "erc20token"
)

type Community struct {
	es      *ethrequest.EthService
	bs      *ethrequest.EthService
	ps      *ethrequest.EthService
	key     *ecdsa.PrivateKey
	address common.Address
	Chain   cw.ChainConfig

	EntryPoint common.Address
	Gateway    *gateway.Gateway

	paddr     common.Address
	Paymaster *paymaster.Paymaster
	ptype     PaymasterType

	afaddr         common.Address
	AccountFactory *simpleaccountfactory.Simpleaccountfactory

	grfaddr          common.Address
	GratitudeFactory *grfactory.Grfactory

	prfaddr        common.Address
	ProfileFactory *profactory.Profactory

	regaddr     common.Address
	RegensToken *regensToken.RegensToken

	deraddr common.Address
	DERC20  *derc20.Derc20

	vu *voucher.VoucherUploader
}

func (c *Community) ExportAddress() CommunityAddress {
	return CommunityAddress{
		Gateway:          c.EntryPoint,
		Paymaster:        c.paddr,
		AccountFactory:   c.afaddr,
		GratitudeFactory: c.grfaddr,
		ProfileFactory:   c.prfaddr,
		RegensToken:      c.regaddr,
		DERC20:           c.deraddr,
		Chain:            c.Chain,
	}
}

// New instantiates a community struct using the provided addresses for the contracts
func New(baseUrl string, es, bs, ps *ethrequest.EthService, ptype string, key *ecdsa.PrivateKey, address common.Address, addr CommunityAddress) (*Community, error) {
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
	acc, err := simpleaccountfactory.NewSimpleaccountfactory(addr.AccountFactory, es.Client())
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

	reg, err := regensToken.NewRegensToken(addr.RegensToken, es.Client())
	if err != nil {
		return nil, err
	}

	der, err := derc20.NewDerc20(addr.DERC20, es.Client())
	if err != nil {
		return nil, err
	}

	vu := voucher.NewVoucherUploader(baseUrl)

	return &Community{
		es:               es,
		bs:               bs,
		ps:               ps,
		key:              key,
		address:          address,
		Chain:            addr.Chain,
		EntryPoint:       addr.Gateway,
		Gateway:          g,
		paddr:            addr.Paymaster,
		Paymaster:        p,
		ptype:            PaymasterType(ptype),
		afaddr:           addr.AccountFactory,
		AccountFactory:   acc,
		grfaddr:          addr.GratitudeFactory,
		GratitudeFactory: gr,
		prfaddr:          addr.ProfileFactory,
		ProfileFactory:   pro,
		regaddr:          addr.RegensToken,
		RegensToken:      reg,
		deraddr:          addr.DERC20,
		DERC20:           der,
		vu:               vu,
	}, nil
}

// Deploy instantiates a community struct and deploys the contracts
func Deploy(baseUrl string, es, bs, ps *ethrequest.EthService, ptype string, key *ecdsa.PrivateKey, address common.Address, chain cw.ChainConfig) (*Community, error) {
	vu := voucher.NewVoucherUploader(baseUrl)

	c := &Community{
		es:      es,
		bs:      bs,
		ps:      ps,
		key:     key,
		address: address,
		Chain:   chain,
		vu:      vu,
		ptype:   PaymasterType(ptype),
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

	err = c.FundPaymaster(big.NewInt(int64(wei.EthToWei(0.1))))
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

	// deploy regens token contract
	err = c.DeployRegensToken()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewTransactor returns a new transactor for the community
func (c *Community) NewTransactor() (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(c.key, big.NewInt(int64(c.Chain.ChainID)))
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
	addr, _, acc, err := simpleaccountfactory.DeploySimpleaccountfactory(auth, c.es.Client(), c.EntryPoint)
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

	_, err = c.AccountFactory.CreateAccount(auth, owner, big.NewInt(int64(0)))
	if err != nil {
		return nil, err
	}

	addr, err := c.AccountFactory.GetAddress(&bind.CallOpts{}, owner, big.NewInt(int64(0)))
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

// DeployRegensToken deploys the regens token contract
func (c *Community) DeployRegensToken() error {
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

	admin1 := common.HexToAddress("0x664ce0F7785E4bA5Ff422C77314eF982F193BeF5")
	admin2 := common.HexToAddress("0xf11704511975cC5908f6dBd89Be922f5C86c1055")

	// deploy profile factory contract
	addr, _, reg, err := regensToken.DeployRegensToken(auth, c.es.Client(), []common.Address{admin1, admin2}, "Regen Treasury", "REG")
	if err != nil {
		return err
	}

	c.regaddr = addr
	c.RegensToken = reg

	return nil
}

// RegensMintCustomVoucher mints a custom voucher for the provided owner
func (c *Community) RegensMintCustomVoucher(owner common.Address, amount *big.Int, minter, name, description string) (*voucher.VoucherMetaData, error) {
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

	// simple id generation
	id := big.NewInt(10000 + int64(nonce))

	// upload the image to ipfs
	resp, err := c.vu.UploadImage(c.regaddr.Hex(), id.String(), minter, owner.Hex(), name, description)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(resp.MetaData)
	if err != nil {
		return nil, err
	}

	_, err = c.RegensToken.MintCustomVoucher(auth, owner, id, amount, b, resp.MetaDataCID)
	if err != nil {
		return nil, err
	}

	return &resp.MetaData, nil
}

// RegensFilterTransferSingle filters the single transfer events for the provided owner, to and/or from
func (c *Community) RegensFilterTransferSingle(owners, froms, tos []common.Address) ([]*regensToken.RegensTokenTransferSingle, error) {

	it, err := c.RegensToken.FilterTransferSingle(&bind.FilterOpts{}, owners, froms, tos)
	if err != nil {
		return nil, err
	}

	t := []*regensToken.RegensTokenTransferSingle{}

	for it.Next() {
		print(it.Event.Id.String())
		t = append(t, it.Event)
	}

	return t, nil
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

// GetDERC20Balance returns the balance of the provided owner
func (c *Community) GetDERC20Balance(owner common.Address) (*big.Int, error) {
	b, err := c.DERC20.BalanceOf(&bind.CallOpts{}, owner)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetPaymasterData returns the paymaster data for the provided userop
func (c *Community) GetPaymasterData(sender common.Address, userop []byte) (*ethrequest.SponsorOp, error) {
	sop, err := c.ps.SponsorUserOp(userop, c.EntryPoint.Hex(), string(c.ptype))
	if err != nil {
		println(err.Error())
		return nil, err
	}

	return sop, nil
}

// SubmitOp submits an operation to the gateway for processing
func (c *Community) SubmitOp(sender common.Address, userop []byte) error {
	err := c.bs.SendUserOp(userop, c.EntryPoint.Hex())
	if err != nil {
		return err
	}

	return nil
}

// setDefaultParameters sets the nonce, value and gas limit for a default contract transaction
func setDefaultParameters(auth *bind.TransactOpts, nonce uint64) {
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(0)
}
