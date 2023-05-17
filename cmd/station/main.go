package main

import (
	"context"
	"flag"
	"log"

	"github.com/daobrussels/cw/pkg/config"
	"github.com/daobrussels/cw/pkg/router"
)

func main() {
	ctx := context.Background()

	log.Default().Println("station starting up...")

	env := flag.String(
		"env",
		"",
		"specify whether to use a dot env file or not",
	)
	port := flag.Int("port", 3000, "specify the port to use")
	_ = flag.String("chainUrl", "http://localhost:8545", "specify the url to use")
	_ = flag.Int("chainId", 80001, "specify the chain id to use")
	flag.Parse()

	conf, err := config.NewConfig(ctx, "chain.json", *env)
	if err != nil {
		log.Fatal(err)
	}

	err = router.NewServer(conf).Start(*port)
	if err != nil {
		log.Fatal(err)
	}

	log.Default().Println("station shutting down...")
}
