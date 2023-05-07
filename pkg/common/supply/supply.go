package supply

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Supply struct {
	PrivateKey    *ecdsa.PrivateKey
	PrivateHexKey string
	PubHexKey     string
	Address       string
}

func New(hexkey string) (*Supply, error) {
	if hexkey == "" {
		return nil, errors.New("SUPPLY_WALLET_KEY is not set")
	}

	privateKey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return nil, err
	}

	compressed := crypto.CompressPubkey(&privateKey.PublicKey)
	if compressed == nil {
		return nil, errors.New("unable to compress public key")
	}

	return &Supply{
		PrivateKey:    privateKey,
		PrivateHexKey: hexkey,
		PubHexKey:     common.Bytes2Hex(compressed),
		Address:       crypto.PubkeyToAddress(privateKey.PublicKey).Hex(),
	}, nil
}
