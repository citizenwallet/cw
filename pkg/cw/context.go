package cw

import "context"

const (
	// SignatureHeader is the header that contains the signature of the request
	SignatureHeader = "X-Signature"
	// PubKeyHeader is used for encrypting the response
	PubKeyHeader = "X-PubKey"
)

type ContextKey string

const (
	ContextKeyPubKey ContextKey = PubKeyHeader
)

// get pub key from context if exists
func GetPubKeyFromContext(ctx context.Context) (string, bool) {
	pubKey, ok := ctx.Value(ContextKeyPubKey).(string)
	return pubKey, ok
}
