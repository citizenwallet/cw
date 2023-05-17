package tests

import (
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/community"
	"github.com/daobrussels/cw/pkg/config"
)

const ()

func TestGateway(t *testing.T) {
	t.Run("test community deploy", func(t *testing.T) {
		ctx := context.Background()

		conf, err := config.NewConfig(ctx, "./data/chain.json", "./data/.test.env")
		if err != nil {
			log.Fatal(err)
		}

		s, err := supply.New(conf.SupplyWalletKey)
		if err != nil {
			log.Fatal(err)
		}

		es, err := ethrequest.NewEthService(conf.Chain.RPC[0])
		if err != nil {
			log.Fatal(err)
		}

		gaddr, err := community.DeployGateway(es, s.PrivateKey, s.Address, big.NewInt(int64(conf.Chain.ChainID)))
		if err != nil {
			log.Fatal(err)
		}

		println("Gateway address:")
		println(gaddr.Hex())

		paddr, err := community.DeployPaymaster(es, s.PrivateKey, s.Address, big.NewInt(int64(conf.Chain.ChainID)), *gaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Paymaster address:")
		println(paddr.Hex())

		afaddr, err := community.DeployAccountFactory(es, s.PrivateKey, s.Address, big.NewInt(int64(conf.Chain.ChainID)), *gaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("Account Factory address:")
		println(afaddr.Hex())

		appaddr, err := community.DeployAppFactory(es, s.PrivateKey, s.Address, big.NewInt(int64(conf.Chain.ChainID)), *gaddr)
		if err != nil {
			log.Fatal(err)
		}

		println("App Factory address:")
		println(appaddr.Hex())
	})
}
