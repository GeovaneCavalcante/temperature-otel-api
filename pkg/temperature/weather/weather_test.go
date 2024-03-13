package weather

import (
	"context"
	"errors"
	"testing"

	"github.com/GeovaneCavalcante/temperatura-cep/pkg/requester"
	mock_requester "github.com/GeovaneCavalcante/temperatura-cep/pkg/requester/mock"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/temperature"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type WeatherFetcherTestSuite struct {
	suite.Suite
	Url             string
	Key             string
	Response        string
	Requester       *mock_requester.MockSender
	TemperatureInfo *temperature.Info
}

func (suite *WeatherFetcherTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.Url = "https://teste/v1/current.json"
	suite.Key = "123456"
	suite.Requester = mock_requester.NewMockSender(ctrl)
	suite.Response = `{"current":{"temp_c":20,"temp_f":60}}`
	suite.TemperatureInfo = &temperature.Info{Celsius: 20, Fahrenheit: 60}
	suite.TemperatureInfo.Kelvin = suite.TemperatureInfo.Celsius + 273.15
}

func (suite *WeatherFetcherTestSuite) TestGetByCity() {
	suite.Run("should return temperature info with successfully", func() {

		suite.Requester.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Return(requester.Response{StatusCode: 200, Body: []byte(suite.Response)}, nil)

		w := New(suite.Url, suite.Key, suite.Requester)

		ctx := context.Background()
		tempInfo, err := w.GetByCity(ctx, "Imperatriz")

		suite.Nil(err)
		suite.NotNil(tempInfo)
		suite.Equal(suite.TemperatureInfo, tempInfo)
	})

	suite.Run("should return error when request return error", func() {
		suite.Requester.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Return(requester.Response{}, errors.New("fail to execute http request for city Imperatriz"))

		w := New(suite.Url, suite.Key, suite.Requester)

		ctx := context.Background()
		tempInfo, err := w.GetByCity(ctx, "Imperatriz")

		suite.Nil(tempInfo)
		suite.NotNil(err)
	})

	suite.Run("should return error when request return status code different of 200", func() {
		suite.Requester.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Return(requester.Response{StatusCode: 500}, nil)

		w := New(suite.Url, suite.Key, suite.Requester)

		ctx := context.Background()
		tempInfo, err := w.GetByCity(ctx, "Imperatriz")

		suite.Nil(tempInfo)
		suite.NotNil(err)
	})
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(WeatherFetcherTestSuite))
}
