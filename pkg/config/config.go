package config

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	// ...
	Port               string `env:"PORT,default=3000"`
	PaymentProviderKey string `env:"PAYMENT_PROVIDER_KEY,required"`
	SupplyWalletKey    string `env:"SUPPLY_WALLßET_KEY,required"`
}

func NewConfig(ctx context.Context, fromFile bool) (*Config, error) {
	if fromFile {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	var conf *Config

	err := envconfig.Process(ctx, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
