package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GeovaneCavalcante/temperatura-cep/internal/entity"
	mock_entity "github.com/GeovaneCavalcante/temperatura-cep/internal/usecase/temperature/mock"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/temperature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type WebTemperatureHandlerTestSuite struct {
	suite.Suite
	FindTemperatureMock *mock_entity.MockFindTemperatureUseCase
}

func (suite *WebTemperatureHandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.FindTemperatureMock = mock_entity.NewMockFindTemperatureUseCase(ctrl)
}

func (suite *WebTemperatureHandlerTestSuite) TestTemperatureByCepHandler() {
	suite.Run("should return temperature info with successfully", func() {
		suite.FindTemperatureMock.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(&temperature.Info{Celsius: 20, Kelvin: 40, Fahrenheit: 60}, nil)

		h := NewWebTemperatureHandler(suite.FindTemperatureMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.TemperatureByCepHandler)
		req, err := http.NewRequest("GET", "/temperature?cep=65909001", nil)
		handler.ServeHTTP(rr, req)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), http.StatusOK, rr.Code)
		assert.Equal(suite.T(), "{\"temp_C\":20,\"temp_F\":60,\"temp_K\":40}\n", rr.Body.String())
	})

	suite.Run("should return erro invalid method", func() {
		h := NewWebTemperatureHandler(suite.FindTemperatureMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.TemperatureByCepHandler)
		req, err := http.NewRequest("POST", "/temperature?cep=65909001", nil)
		handler.ServeHTTP(rr, req)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), http.StatusMethodNotAllowed, rr.Code)
		assert.Equal(suite.T(), "", rr.Body.String())
	})

	suite.Run("should return erro invalid zipcode", func() {
		h := NewWebTemperatureHandler(suite.FindTemperatureMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.TemperatureByCepHandler)
		req, err := http.NewRequest("GET", "/temperature?cep=65909001A", nil)
		handler.ServeHTTP(rr, req)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), http.StatusUnprocessableEntity, rr.Code)
		assert.Equal(suite.T(), "invalid zipcode", rr.Body.String())
	})

	suite.Run("should return erro zipcode not found", func() {
		suite.FindTemperatureMock.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, entity.ErrZipCodeNotFound)

		h := NewWebTemperatureHandler(suite.FindTemperatureMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.TemperatureByCepHandler)
		req, err := http.NewRequest("GET", "/temperature?cep=65909001", nil)
		handler.ServeHTTP(rr, req)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), http.StatusNotFound, rr.Code)
		assert.Equal(suite.T(), "can not found zipcode", rr.Body.String())
	})

	suite.Run("should return erro when usecase return error", func() {
		suite.FindTemperatureMock.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		h := NewWebTemperatureHandler(suite.FindTemperatureMock)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.TemperatureByCepHandler)
		req, err := http.NewRequest("GET", "/temperature?cep=65909001", nil)
		handler.ServeHTTP(rr, req)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)
		assert.Equal(suite.T(), "internal server error", rr.Body.String())
	})

}
func TestSuite(t *testing.T) {
	suite.Run(t, new(WebTemperatureHandlerTestSuite))
}
