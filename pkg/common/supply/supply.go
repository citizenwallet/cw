package supply

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

type Supply struct {
	PrivateHexKey string
}

var PrivateHexKey string
var PublicHexKey string

func init() {
	hexkey := os.Getenv("SUPPLY_WALLET_KEY")
	if hexkey == "" {
		log.Fatal("missing supply key")
	}

	PrivateHexKey = hexkey

	privateKey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		log.Fatal(err)
	}

	compressed := crypto.CompressPubkey(&privateKey.PublicKey)
	if compressed == nil {
		log.Fatal("invalid public key")
	}

	PublicHexKey = string(compressed)
}
