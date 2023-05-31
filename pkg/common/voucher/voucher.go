package voucher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/daobrussels/cw/pkg/common/request"
)

const (
	UploadURLTemplate = "/api/mint?contract_chain=polygon&contract_address=%s&id=%s&minter_name=%s&minter_address=%s&name=%s&description=%s"
)

type VoucherUploader struct {
	baseURL string
}

func NewVoucherUploader(baseURL string) *VoucherUploader {
	return &VoucherUploader{
		baseURL: baseURL,
	}
}

type IPFSUploadResponse struct {
	ContractAddress string          `json:"contract_address"`
	MetaDataCID     string          `json:"metadata_cid"`
	MetaData        VoucherMetaData `json:"metadata"`
}

type VoucherMetaData struct {
	Name            string `json:"name"`
	Desc            string `json:"description"`
	MinterName      string `json:"minter_name"`
	MinterAddress   string `json:"minter_address"`
	ContractChain   string `json:"contract_chain"`
	ContractAddress string `json:"contract_address"`
	MintingDate     string `json:"minting_date"`
	Image           string `json:"image"`
}

func (v *VoucherUploader) UploadImage(caddr, id, minter, maddr, name, description string) (*IPFSUploadResponse, error) {
	url := fmt.Sprintf(v.baseURL+UploadURLTemplate, caddr, id, url.QueryEscape(minter), maddr, url.QueryEscape(name), url.QueryEscape(description))

	r, err := request.HttpRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	if r.StatusCode == http.StatusForbidden {
		return nil, errors.New("auth failed")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.New("failed to upload image")
	}

	var resp IPFSUploadResponse
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
