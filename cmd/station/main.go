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
	"github.com/daobrussels/cw/pkg/router"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	ctx := context.Background()

	log.Default().Println("station starting up...")

	env := flag.String(
		"env",
		".env",
		"specify path to env",
	)

	port := flag.Int(
		"port",
		3000,
		"specify port to listen on",
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

	es, err := ethrequest.NewEthService(addr.Chain.RPC[0])
	if err != nil {
		log.Fatal(err)
	}
	defer es.Close()

	conf, err := config.NewConfigWChain(ctx, *env, addr.Chain)
	if err != nil {
		log.Default().Println(fmt.Sprintf("invalid or missing chain config file at %s", *path))
		log.Fatal(err)
	}

	s, err := supply.New(conf.SupplyWalletKey)
	if err != nil {
		log.Fatal(err)
	}

	c, err := community.New(es, s.PrivateKey, common.HexToAddress(s.Address), addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Default().Println("serving...")

	err = router.NewServer(s, es, c).Start(*port)
	if err != nil {
		log.Fatal(err)
	}

	log.Default().Println("station shutting down...")
}
