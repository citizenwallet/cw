package tests

import (
	"testing"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/common/transaction"
	"github.com/daobrussels/cw/pkg/cw"
)

const (
	txrpc        = "http://localhost:8545"
	txprivhexkey = "baa1ad2bed05fd3a9b2c8b5feeb98b464454ca2586e16d58768df8dcf03d4b91"
	// txaddress          = "0xdcb53b8667dbef2b8cfcba912718bf3ba3173913"
	txreceivingAddress = "0xe13b2276bb63fde321719bbf6dca9a70fc40efcc"
)

func TestTransaction(t *testing.T) {
	t.Run("test transaction", func(t *testing.T) {
		supply, err := supply.New(txprivhexkey)
		if err != nil {
			t.Fatal(err)
		}

		ethservice, err := ethrequest.NewEthService(txrpc)
		if err != nil {
			t.Fatal(err)
		}
		defer ethservice.Close()

		chain, err := cw.GetChain("data/chain.json")
		if err != nil {
			t.Fatal(err)
		}

		s := transaction.New(chain, supply, ethservice)

		err = s.Send(txreceivingAddress, 1000000000000000000)
		if err != nil {
			t.Fatal(err)
		}
	})
}
