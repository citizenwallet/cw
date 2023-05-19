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
	"github.com/ethereum/go-ethereum/common"
)

const (
	nobalancehexaddr  = "0x664ce0F7785E4bA5Ff422C77314eF982F193BeF5"
	nobalancehexaddr2 = "0xb9721dab9c7eca07d836afdc6244e22be2ba91ea"
)

func TestCommunity(t *testing.T) {
	var gaddr, afaddr, pfaddr, graddr *common.Address

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

		gaddr, err = community.DeployGateway(es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)))
		if err != nil {
			log.Fatal(err)
		}

		println("Gateway address:")
		println(gaddr.Hex())

		paddr, pay, err := community.DeployPaymaster(es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)), *gaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Paymaster address:")
		println(paddr.Hex())

		amount := big.NewInt(int64(wei.EthToWei(1)))
		err = community.FundPaymaster(pay, es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)), amount)
		if err != nil {
			log.Fatal(err)
		}

		println("Paymaster funded:")
		println(amount.String())

		afaddr, err = community.DeployAccountFactory(es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)), *gaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Account Factory address:")
		println(afaddr.Hex())

		pfaddr, err = community.DeployProfileFactory(es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)), *gaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Profile Factory address:")
		println(pfaddr.Hex())

		graddr, err = community.DeployGratitudeFactory(es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)), *gaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Gratitude Factory address:")
		println(graddr.Hex())
	})

	t.Run("test account deploy on community", func(t *testing.T) {
		owner := common.HexToAddress(nobalancehexaddr)

		// create an account
		accaddr, err := community.CreateAccount(es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)), owner, *afaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Account address:")
		println(accaddr.Hex())

		grtaddr, err := community.CreateGratitudeApp(es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)), *accaddr, *graddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Gratitude address:")
		println(grtaddr.Hex())

		// create a profile for the corresponding account
		profile, err := community.CreateProfile(es, s.PrivateKey, maddress, big.NewInt(int64(conf.Chain.ChainID)), *accaddr, *pfaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Profile address:")
		println(profile.Hex())
	})
}
