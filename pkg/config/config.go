package config

import (
	"context"

	"github.com/daobrussels/cw/pkg/cw"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	// ...
	PaymentProviderKey string `env:"PAYMENT_PROVIDER_KEY,required"`
	SupplyWalletKey    string `env:"SUPPLY_WALLET_KEY,required"`
	Chain              cw.ChainConfig
}

func NewConfig(ctx context.Context, fromFile bool) (*Config, error) {
	if fromFile {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	conf := &Config{}

	err := envconfig.Process(ctx, conf)
	if err != nil {
		return nil, err
	}

	chain, err := cw.GetChain("chain.json")
	if err != nil {
		return nil, err
	}

	conf.Chain = *chain

	return conf, nil
}
