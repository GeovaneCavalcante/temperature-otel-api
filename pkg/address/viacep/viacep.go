package viacep

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GeovaneCavalcante/temperatura-cep/pkg/address"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/logger"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/requester"
)

type Response struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        bool   `json:"erro"`
}

type ViaCepFetcher struct {
	Url       string
	Requester requester.Sender
}

func New(url string, req requester.Sender) *ViaCepFetcher {
	return &ViaCepFetcher{
		Url:       url,
		Requester: req,
	}
}

func (v *ViaCepFetcher) GetByZipCode(ctx context.Context, zipCode string) (*address.Info, error) {

	logger.Info("[ViaCepFetcher] starting api for zipcode: " + zipCode)

	cfg := requester.Configuration{
		Url:        v.Url + "/" + zipCode + "/json/",
		Method:     "get",
		ContetType: "application/json",
	}

	res, err := v.Requester.Send(ctx, cfg)
	if err != nil {
		logger.Error("[ViaCepFetcher] fail to execute http request for zipcode "+zipCode, err)
		return nil, fmt.Errorf("fail to execute http request for zipcode %s: %w", zipCode, err)
	}

	var zipCodeInfo Response

	err = json.Unmarshal(res.Body, &zipCodeInfo)
	if err != nil {
		logger.Error("[ViaCepFetcher] fail to unmarshal response body for zipcode "+zipCode, err)
		return nil, fmt.Errorf("fail to unmarshal response body for zipcode %s: %w", zipCode, err)
	}

	if zipCodeInfo.Erro {
		logger.Error("[ViaCepFetcher] zipcode not found: "+zipCode, nil)
		return nil, fmt.Errorf("zipcode not found: %s", zipCode)
	}
	AddressInfo := toAddressInfo(zipCodeInfo)

	logger.Info("[ViaCepFetcher] finishing api for zipcode " + zipCode)

	return AddressInfo, nil
}

func toAddressInfo(zipCodeInfo Response) *address.Info {
	return &address.Info{
		City: zipCodeInfo.Localidade,
	}
}
