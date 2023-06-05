package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/daobrussels/cw/pkg/common/ethrequest"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/community"
	"github.com/daobrussels/cw/pkg/config"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	ctx := context.Background()

	log.Default().Println("preparing...")

	env := flag.String(
		"env",
		".env",
		"specify path to env",
	)

	gwei := flag.Int64(
		"gwei",
		100000000000000000,
		"specify amount of gwei to fund",
	)

	path := flag.String(
		"c",
		"./config/community/test.community.json",
		"specify path to a *.community.json file",
	)

	flag.Parse()

	b, err := os.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}

	var addr community.CommunityAddress
	err = json.Unmarshal(b, &addr)
	if err != nil {
		log.Fatal(err)
	}

	conf, err := config.NewConfigWChain(ctx, *env, addr.Chain)
	if err != nil {
		log.Default().Println(fmt.Sprintf("invalid or missing chain config file at %s", *path))
		log.Fatal(err)
	}

	es, err := ethrequest.NewEthService(addr.Chain.RPC[0])
	if err != nil {
		log.Fatal(err)
	}
	defer es.Close()

	bs, err := ethrequest.NewEthService(conf.RPCUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer bs.Close()

	ps, err := ethrequest.NewEthService(conf.PaymasterRPCUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer ps.Close()

	s, err := supply.New(conf.SupplyWalletKey)
	if err != nil {
		log.Fatal(err)
	}

	c, err := community.New(conf.RegensIPFSUploader, es, bs, ps, conf.PaymasterType, s.PrivateKey, common.HexToAddress(s.Address), addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Default().Println("funding...")

	err = c.FundPaymaster(big.NewInt(*gwei))
	if err != nil {
		log.Fatal(err)
	}

	log.Default().Println("funded...")
}
