package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/daobrussels/cw/pkg/common/request"
	"github.com/daobrussels/cw/pkg/common/supply"
	"github.com/daobrussels/cw/pkg/cw"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type ResponseType string

const (
	ResponseTypeObject ResponseType = "object"
	ResponseTypeArray  ResponseType = "array"
)

type Response struct {
	ResponseType ResponseType `json:"response_type"`
	Secure       string       `json:"secure,omitempty"`
	Object       any          `json:"object,omitempty"`
	Objects      any          `json:"objects,omitempty"`
}

func Body(w http.ResponseWriter, body any) error {

	w.Header().Add("Content-Type", "application/json")

	b, err := json.Marshal(&Response{
		ResponseType: ResponseTypeObject,
		Object:       body,
	})
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func EncryptedBody(w http.ResponseWriter, ctx context.Context, body any) error {

	pubhexkey, ok := cw.GetPubKeyFromContext(ctx)
	if !ok {
		return errors.New("unable to parse public key from context")
	}

	pubkey, err := crypto.DecompressPubkey(common.Hex2Bytes(pubhexkey))
	if err != nil {
		return err
	}

	address := crypto.PubkeyToAddress(*pubkey)

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := request.New(address.Hex(), b)

	sig, err := req.GenerateSignature(supply.PrivateHexKey)
	if err != nil {
		return err
	}

	secure, err := req.Encrypt(pubhexkey)
	if err != nil {
		return err
	}

	bresp, err := json.Marshal(&Response{
		ResponseType: ResponseTypeObject,
		Secure:       secure,
	})
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add(cw.SignatureHeader, sig)
	w.Header().Add(cw.PubKeyHeader, supply.PublicHexKey)
	w.Write(bresp)

	return nil
}

// func BodyMultiple(w http.ResponseWriter, body any) error {

// 	w.Header().Add("ContentType", "application/json")

// 	b, err := json.Marshal(&Response{
// 		ResponseType: ResponseTypeArray,
// 		Objects:      body,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	w.Write(b)

// 	return nil
// }

// func StreamedBody(w http.ResponseWriter, body string) error {
// 	flusher, ok := w.(http.Flusher)
// 	if !ok {
// 		return errors.New("stearming not supported")
// 	}

// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	fmt.Fprintf(w, "%s", body)
// 	flusher.Flush()

// 	return nil
// }
