package request

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
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
		Expiry:  time.Now().Add(10 * time.Second),
		Address: address,
		Data:    data,
	}
}

// Encrypt encrypts the request data using the public key, result is base64 encoded
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

	// base64 encode the encrypted data
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts the base64 encoded request data using a private key
func Decrypt(hexkey string, req string) (*Request, error) {
	b, err := base64.StdEncoding.DecodeString(req)
	if err != nil {
		return nil, err
	}

	// decode the private key
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

	// has the request expired?
	if time.Now().After(r.Expiry) {
		return false
	}

	// hash the request data
	h := crypto.HashData(crypto.NewKeccakState(), b)

	// decode the signature
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	// recover the public key from the signature
	pubkey, err := crypto.SigToPub(h.Bytes(), sig)
	if err != nil {
		return false
	}

	// derive the address from the public key
	address := crypto.PubkeyToAddress(*pubkey)

	// the address in the request must match the address derived from the signature
	if address.Hex() != string(r.Address) {
		return false
	}

	// compress the public key
	compressed := crypto.CompressPubkey(pubkey)

	// remove the recovery id from the signature
	cleanSig := sig[:len(sig)-1]

	// verify the signature with the derived public key and the hash of the request data
	return crypto.VerifySignature(compressed, h.Bytes(), cleanSig)
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

	// hash the request data
	h := crypto.HashData(crypto.NewKeccakState(), b)

	// sign the hash of the request data
	s, err := crypto.Sign(h.Bytes(), privateKey)

	// base64 encode the signature
	return base64.StdEncoding.EncodeToString(s), nil
}
