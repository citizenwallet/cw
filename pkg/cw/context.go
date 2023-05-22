package cw

import (
	"context"
)

const (
	// SignatureHeader is the header that contains the signature of the request
	SignatureHeader = "X-Signature"
	// PubKeyHeader is used for encrypting the response
	PubKeyHeader = "X-PubKey"
	// AddressHeader is the header that contains the address of the sender
	AddressHeader = "X-Address"
)

type ContextKey string

const (
	ContextKeyPubKey  ContextKey = PubKeyHeader
	ContextKeyAddress ContextKey = AddressHeader
)

// get pub key from context if exists
func GetPubKeyFromContext(ctx context.Context) (string, bool) {
	pubKey, ok := ctx.Value(ContextKeyPubKey).(string)
	return pubKey, ok
}

// get address from context if exists
func GetAddressFromContext(ctx context.Context) (string, bool) {
	addr, ok := ctx.Value(ContextKeyAddress).(string)
	return addr, ok
}
