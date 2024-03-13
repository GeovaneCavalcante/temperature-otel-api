package temperature

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GeovaneCavalcante/temperatura-cep/internal/entity"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/logger"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/requester"
)

type ProxyTemperatureUseCase struct {
	Url       string
	Requester requester.Sender
}

func NewProxyTemperatureUseCase(url string, requester requester.Sender) *ProxyTemperatureUseCase {
	return &ProxyTemperatureUseCase{
		Url:       url,
		Requester: requester,
	}
}

func (s *ProxyTemperatureUseCase) Execute(ctx context.Context, zipCode string) (*entity.TemperatureInfo, error) {

	logger.Info("[ProxyTemperatureUseCase] starting usecase for zipcode: " + zipCode)

	cfg := requester.Configuration{
		Url:        s.Url + "?cep=" + zipCode,
		Method:     "get",
		ContetType: "application/json",
	}

	res, err := s.Requester.Send(ctx, cfg)

	if err != nil {
		logger.Error("[ProxyTemperatureUseCase] fail to execute http request for zipcode "+zipCode, err)
		return nil, fmt.Errorf("fail to execute http request for zipcode %s", zipCode)
	}
	var tempInfo entity.TemperatureInfo

	err = json.Unmarshal(res.Body, &tempInfo)
	if err != nil {
		logger.Error("[ProxyTemperatureUseCase] fail to unmarshal response body for zipcode "+zipCode, err)
		return nil, err
	}

	logger.Info("[ProxyTemperatureUseCase] finishing usecase for zipcode: " + zipCode + " with success")

	return &tempInfo, nil
}
