package main

import (
	"context"
	"flag"
	"log"

	"github.com/daobrussels/cw/pkg/config"
)

func main() {
	ctx := context.Background()

	log.Default().Println("events starting up...")

	env := flag.Bool(
		"env",
		false,
		"specify whether to use a dot env file or not",
	)
	_ = flag.String("url", "http://localhost:8545", "specify the url to use")
	flag.Parse()

	_, err := config.NewConfig(ctx, *env)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: implement the events server

	log.Default().Println("events shutting down...")
}
