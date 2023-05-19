package tests

import (
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/common/wei"
	"github.com/daobrussels/cw/pkg/community"
	"github.com/daobrussels/cw/pkg/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const (
	nobalancehexaddr  = "0x664ce0F7785E4bA5Ff422C77314eF982F193BeF5"
	nobalancehexaddr2 = "0xb9721dab9c7eca07d836afdc6244e22be2ba91ea"
)

func TestCommunity(t *testing.T) {
	var entryPoint common.Address
	var c *community.Community

	ctx := context.Background()

	conf, err := config.NewConfig(ctx, "./data/chain.json", "./data/.test.env")
	if err != nil {
		log.Fatal(err)
	}

	s, err := supply.New(conf.SupplyWalletKey)
	if err != nil {
		log.Fatal(err)
	}

	maddress := common.HexToAddress(s.Address)

	es, err := ethrequest.NewEthService(conf.Chain.RPC[0])
	if err != nil {
		log.Fatal(err)
	}

	t.Run("test community deploy", func(t *testing.T) {

		// deploy community
		c, err = community.Deploy(es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)))
		if err != nil {
			log.Fatal(err)
		}

		entryPoint = c.EntryPoint

		println("Community Entrypoint:")
		println(c.EntryPoint.Hex())

		amount := big.NewInt(int64(wei.EthToWei(1)))
		err = c.FundPaymaster(amount)
		if err != nil {
			log.Fatal(err)
		}

		println("Paymaster funded:")
		println(amount.String())
	})

	t.Run("test account deploy on community", func(t *testing.T) {
		owner := common.HexToAddress(nobalancehexaddr)

		// create an account
		accaddr, err := c.CreateAccount(owner)
		if err != nil {
			log.Fatal(err)
		}

		println("Account address:")
		println(accaddr.Hex())

		grtaddr, err := c.CreateGratitudeApp(*accaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Gratitude address:")
		println(grtaddr.Hex())

		// create a profile for the corresponding account
		profile, err := c.CreateProfile(*accaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Profile address:")
		println(profile.Hex())

		// get account
		acc, err := c.GetAccount(*accaddr)
		if err != nil {
			log.Fatal(err)
		}

		acccep, err := acc.EntryPoint(&bind.CallOpts{})
		if err != nil {
			log.Fatal(err)
		}

		println("Account entrypoint address:")
		println(acccep.Hex())

		if acccep.Hex() != entryPoint.Hex() {
			log.Fatal("Account entrypoint address is not the same as community entrypoint address")
		}

		// get account
		pr, err := c.GetProfile(*accaddr)
		if err != nil {
			log.Fatal(err)
		}

		pep, err := pr.EntryPoint(&bind.CallOpts{})
		if err != nil {
			log.Fatal(err)
		}

		println("Profile:")
		println(pep.Hex())

		if pep.Hex() != entryPoint.Hex() {
			log.Fatal("Profile entrypoint address is not the same as community entrypoint address")
		}
	})
}
