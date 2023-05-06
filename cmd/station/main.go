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

	env := flag.Bool(
		"env",
		false,
		"specify whether to use a dot env file or not",
	)
	port := flag.Int("port", 3000, "specify the port to use")
	_ = flag.String("url", "http://localhost:8545", "specify the url to use")
	flag.Parse()

	conf, err := config.NewConfig(ctx, *env)
	if err != nil {
		log.Fatal(err)
	}

	err = router.NewServer(conf).Start(*port)
	if err != nil {
		log.Fatal(err)
	}

	log.Default().Println("station shutting down...")
}
