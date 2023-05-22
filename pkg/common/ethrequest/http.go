package ethrequest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ErrRequest error

const vsn string = "2.0"

var (
	ErrBadRequest        ErrRequest = errors.New("invalid request")
	ErrRateLimitExceeded ErrRequest = errors.New("rate limit exceeded")
)

type jsonrpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

type RawService struct {
	url string
}

func NewRawService(url string) *RawService {
	return &RawService{url}
}

func (r *RawService) Post(method string, body []any) ([]byte, error) {
	msg := jsonrpcMessage{Version: vsn, Method: method}

	json_data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	msg.Params = json_data

	msg_data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, r.url, bytes.NewBuffer(msg_data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: time.Second * 60}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 400:
		return nil, ErrBadRequest
	case 429:
		return nil, ErrRateLimitExceeded
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, err := io.ReadAll(resp.Body)
		if err == nil {
			println(string(b))
		}

		return nil, fmt.Errorf("req: failed status code %d", resp.StatusCode)
	}

	type respbody struct {
		Result json.RawMessage `json:"result,omitempty"`
	}

	var rbody respbody
	err = json.NewDecoder(resp.Body).Decode(&rbody)
	if err != nil {
		return nil, err
	}

	return rbody.Result, nil
}
