package cw

import (
	"encoding/json"
	"os"
)

type ChainFeature struct {
	Name string `json:"name"`
}

type ChainNativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

type ChainENS struct {
	Registry string `json:"registry"`
}

type ChainExplorer struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Standard string `json:"standard"`
}

// ChainConfig is the configuration for a chain
type ChainConfig struct {
	Name           string              `json:"name"`
	Chain          string              `json:"chain"`
	Icon           string              `json:"icon"`
	RPC            []string            `json:"rpc"`
	Features       []ChainFeature      `json:"features"`
	Faucets        []string            `json:"faucets"`
	NativeCurrency ChainNativeCurrency `json:"nativeCurrency"`
	InfoURL        string              `json:"infoURL"`
	ShortName      string              `json:"shortName"`
	ChainID        int                 `json:"chainID"`
	NetworkID      int                 `json:"networkID"`
	Slip44         int                 `json:"slip44"`
	ENS            ChainENS            `json:"ens"`
	Explorers      []ChainExplorer     `json:"explorers"`
}

// GetChain returns the chain config for the local chain.json file
func GetChain(path string) (*ChainConfig, error) {
	// read the chain.json file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	println(f.Name())

	// decode the chain.json file
	chain := &ChainConfig{}
	err = json.NewDecoder(f).Decode(chain)
	if err != nil {
		return nil, err
	}

	return chain, nil
}
