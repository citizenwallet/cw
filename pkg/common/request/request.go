package request

import (
	"encoding/json"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	bitcoin_ecies "github.com/gitzhou/bitcoin-ecies"
)

const (
	hexPadding = "0x"
)

type Request struct {
	Version int       `json:"version,required"` // version of the request
	Expiry  time.Time `json:"expiry,required"`  // avoid replay attacks
	Address string    `json:"address,required"` // address of the sender, must match the signature
	Data    []byte    `json:"data,required"`    // data to be sent to the server
}

func New(address string, data []byte) *Request {
	return &Request{
		Version: 1,
		Expiry:  time.Now().Add(10 * time.Second).UTC(),
		Address: address,
		Data:    data,
	}
}

// Encrypt encrypts the request data using the public key, result is base64 encoded
func (r *Request) Encrypt(pubhexkey string) (string, error) {
	publicKey, err := secp256k1.ParsePubKey(common.Hex2Bytes(pubhexkey))
	if err != nil {
		return "", err
	}

	// marshal the request to bytes
	msg, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	encrypted, err := bitcoin_ecies.EncryptMessage(string(msg), publicKey.SerializeUncompressed())
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

// Decrypt decrypts the base64 encoded request data using a private key
func Decrypt(hexkey string, req string) (*Request, error) {
	// decrypt the request data
	decrypted, err := bitcoin_ecies.DecryptMessage(req, common.Hex2Bytes(hexkey))
	if err != nil {
		return nil, err
	}

	// unmarshal the request
	r := &Request{}
	err = json.Unmarshal([]byte(decrypted), r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// VerifySignature verifies the provided signature using the public key against a marshalled version of the request
func (r *Request) VerifySignature(signature string) bool {
	// marshal the request to bytes
	b, err := json.Marshal(r)
	if err != nil {
		return false
	}

	// has the request expired?
	if time.Now().After(r.Expiry) {
		return false
	}

	// hash the request data
	h := crypto.Keccak256Hash(b)

	// decode the signature
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false
	}

	// recover the public key from the signature
	pubkey, _, err := ecdsa.RecoverCompact(sig, h.Bytes())
	if err != nil {
		return false
	}

	// derive the address from the public key
	address := common.BytesToAddress(pubkey.SerializeUncompressed())

	// the address in the request must match the address derived from the signature
	if address.Hex() != string(r.Address) {
		return false
	}

	// create ModNScalars from the signature manually
	sr, ss := secp256k1.ModNScalar{}, secp256k1.ModNScalar{}

	// set the byteslices manually from the signature
	sr.SetByteSlice(sig[1:33])
	ss.SetByteSlice(sig[33:65])

	// create a new signature from the ModNScalars
	ns := ecdsa.NewSignature(&sr, &ss)
	if err != nil {
		return false
	}

	// verify the signature
	return ns.Verify(h.Bytes(), pubkey)
}

// GenerateSignature generates a signature for the request using a private key
func (r *Request) GenerateSignature(hexkey string) (string, error) {
	// marshal the request to bytes
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	// hash the request data
	h := crypto.Keccak256Hash(b)

	privateKey := secp256k1.PrivKeyFromBytes(common.Hex2Bytes(hexkey))

	// sign the hash of the request data
	s := ecdsa.SignCompact(privateKey, h.Bytes(), false)
	if s == nil {
		return "", err
	}

	// base64 encode the signature
	return hexutil.Encode(s), nil
}
