package request

import "encoding/json"

type Request struct {
	Version string `json:"version"`
	Address string `json:"address"`
	Data    string `json:"data"`
}

// VerifySignature verifies the provided signature using the public key against a marshalled version of the request
func (r *Request) VerifySignature(signature, pubkey []byte) bool {
	// marshal the request to bytes
	_, err := json.Marshal(r)
	if err != nil {
		return false
	}

	// encrypt the marshalled request using the public key

	return true
}
