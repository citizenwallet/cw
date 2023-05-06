package request

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const (
	hexPadding = "0x"
)

type Request struct {
	Version string `json:"version"`
	Address []byte `json:"address"`
	Data    []byte `json:"data"`
}

func New(address, data []byte) *Request {
	return &Request{
		Version: "1.0",
		Address: address,
		Data:    data,
	}
}

// Encrypt encrypts the request data using the public key
func (r *Request) Encrypt(pubhexkey string) (string, error) {
	publicKey, err := crypto.DecompressPubkey(common.Hex2Bytes(pubhexkey))
	if err != nil {
		return "", err
	}

	// marshal the request to bytes
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	// encrypt the request data
	encrypted, err := ecies.Encrypt(rand.Reader, ecies.ImportECDSAPublic(publicKey), b, nil, nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts the request data using the private key
func Decrypt(hexkey string, req string) (*Request, error) {
	b, err := base64.StdEncoding.DecodeString(req)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return nil, err
	}

	// decrypt the request data
	decrypted, err := ecies.ImportECDSA(privateKey).Decrypt(b, nil, nil)
	if err != nil {
		return nil, err
	}

	// unmarshal the request
	r := &Request{}
	err = json.Unmarshal(decrypted, r)
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

	h := crypto.HashData(crypto.NewKeccakState(), b)

	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	pubkey, err := crypto.SigToPub(h.Bytes(), sig)
	if err != nil {
		return false
	}

	address := crypto.PubkeyToAddress(*pubkey)

	// the address in the request must match the address derived from the signature
	if address.Hex() != string(r.Address) {
		return false
	}

	compressed := crypto.CompressPubkey(pubkey)

	// remove the recovery id from the signature
	return crypto.VerifySignature(compressed, h.Bytes(), sig[:len(sig)-1])
}

// GenerateSignature generates a signature for the request using a private key
func (r *Request) GenerateSignature(hexkey string) (string, error) {
	privateKey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return "", err
	}

	// marshal the request to bytes
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	h := crypto.HashData(crypto.NewKeccakState(), b)

	s, err := crypto.Sign(h.Bytes(), privateKey)

	return base64.StdEncoding.EncodeToString(s), nil
}
