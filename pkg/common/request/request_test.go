package request

import (
	"encoding/json"
	"testing"
)

const (
	privhexkey = "b123284ed609ca4c19a78124567d606f1202630e72784602475f1eb0b3f0a0a2"
	pubhexkey  = "0288cd52ce87d3e674a2383f009e2c956402b99675bc1dc0414a4b78d98dde634b"
	address    = "0x39bC81005a2BEa2122A2F2fd963Db3ac8aDbC518"
)

type TestData struct {
	Hello string `json:"hello"`
}

func TestRequest(t *testing.T) {
	t.Run("test request encryption", func(t *testing.T) {
		b, err := json.Marshal(TestData{Hello: "world"})
		if err != nil {
			t.Fatal(err)
		}

		// test request signature
		req := New(address, b)

		// encrypt the request data
		encrypted, err := req.Encrypt(pubhexkey)
		if err != nil {
			t.Fatal(err)
		}

		// decrypt the request data
		decrypted, err := Decrypt(privhexkey, encrypted)
		if err != nil {
			t.Fatal(err)
		}

		// verify the request data
		if string(decrypted.Data) != string(req.Data) {
			t.Fatal("decrypted data does not match original data")
		}
	})

	t.Run("test request signature", func(t *testing.T) {
		b, err := json.Marshal(TestData{Hello: "world"})
		if err != nil {
			t.Fatal(err)
		}

		req := New(address, b)

		// generate signature
		sig, err := req.GenerateSignature(privhexkey)
		if err != nil {
			t.Fatal(err)
		}

		// verify signature
		if !req.VerifySignature(sig) {
			t.Fatal("signature verification failed")
		}
	})
}
