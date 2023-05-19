package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/community"
	"github.com/daobrussels/cw/pkg/config"
	"github.com/ethereum/go-ethereum/common"
)

const (
	chainTemplate = `
	{
		"name": "Chain Name",
		"chain": "CHAINSYMBOL",
		"icon": "icon_name",
		"rpc": [
			"rpc_url"
		],
		"features": [
			{
				"name": "EIP155"
			},
			{
				"name": "EIP1559"
			}
		],
		"faucets": [],
		"nativeCurrency": {
			"name": "Coin Name",
			"symbol": "Coin Symbol",
			"decimals": 2
		},
		"infoURL": "info_url"",
		"shortName": "cc",
		"chainId": 123456,
		"networkId": 123456,
		"slip44": 60,
		"ens": {
			"registry": "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"
		},
		"explorers": [
			{
				"name": "etherscan",
				"url": "https://etherscan.io",
				"standard": "EIP3091"
			}
		]
	}
	`
)

func main() {
	ctx := context.Background()

	log.Default().Println("deploying...")

	env := flag.String(
		"env",
		".env",
		"specify path to env",
	)

	flag.Parse()

	conf, err := config.NewConfig(ctx, "chain.json", *env)
	if err != nil {
		log.Default().Println("invalid or missing chain.json file")
		log.Default().Println("should be:")
		log.Default().Println(chainTemplate)
		log.Default().Println("")
		log.Default().Println("put at the base of the project")
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

	c, err := community.Deploy(es, s.PrivateKey, maddress, conf.Chain)
	if err != nil {
		log.Fatal(err)
	}

	addr := c.ExportAddress()

	b, err := json.Marshal(addr)
	if err != nil {
		log.Fatal(err)
	}

	// write bytes to file
	err = os.WriteFile(fmt.Sprintf("%s_community.json", addr.Gateway.Hex()), b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Default().Println("community deployed...")
}
