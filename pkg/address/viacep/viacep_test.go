package viacep

import (
	"context"
	"errors"
	"testing"

	"github.com/GeovaneCavalcante/temperatura-cep/pkg/requester"
	mock_requester "github.com/GeovaneCavalcante/temperatura-cep/pkg/requester/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ViaCepTestSuite struct {
	suite.Suite
	Url       string
	Requester *mock_requester.MockSender
	Response  string
}

func (suite *ViaCepTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.Url = "https://teste/v1/current.json"
	suite.Requester = mock_requester.NewMockSender(ctrl)
	suite.Response = `{"cep":"65903-270","logradouro":"Rua Jo\u00e3o Lisboa","complemento":"","bairro":"Bacuri","localidade":"Imperatriz","uf":"MA","ibge":"2105302","gia":"","ddd":"99","siafi":"0921","erro":false}`
}

func (suite *ViaCepTestSuite) TestGetByZipCode() {
	suite.Run("should return address info with successfully", func() {

		suite.Requester.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Return(requester.Response{StatusCode: 200, Body: []byte(suite.Response)}, nil)

		v := New(suite.Url, suite.Requester)

		ctx := context.Background()
		addressInfo, err := v.GetByZipCode(ctx, "65903270")

		suite.Nil(err)
		suite.NotNil(addressInfo)
		suite.Equal("Imperatriz", addressInfo.City)
	})

	suite.Run("should return error when request return error", func() {
		suite.Requester.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Return(requester.Response{}, nil)

		v := New(suite.Url, suite.Requester)
		ctx := context.Background()

		addressInfo, err := v.GetByZipCode(ctx, "65903270")

		suite.Nil(addressInfo)
		suite.NotNil(err)
		suite.Equal("fail to unmarshal response body for zipcode 65903270: unexpected end of JSON input", err.Error())
	})

	suite.Run("should return error when request return error", func() {
		suite.Requester.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Return(requester.Response{}, errors.New("fail to execute http request for zipcode 65903270"))

		v := New(suite.Url, suite.Requester)
		ctx := context.Background()

		addressInfo, err := v.GetByZipCode(ctx, "65903270")

		suite.Nil(addressInfo)
		suite.NotNil(err)
	})

	suite.Run("should return error when request return error", func() {
		suite.Requester.EXPECT().
			Send(gomock.Any(), gomock.Any()).
			Return(requester.Response{StatusCode: 404, Body: []byte(`{"erro":true}`)}, nil)

		v := New(suite.Url, suite.Requester)
		ctx := context.Background()

		addressInfo, err := v.GetByZipCode(ctx, "65903270")

		suite.Nil(addressInfo)
		suite.NotNil(err)
		suite.Equal("zipcode not found: 65903270", err.Error())
	})
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ViaCepTestSuite))
}
