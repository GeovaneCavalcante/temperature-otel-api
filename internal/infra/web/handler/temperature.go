package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GeovaneCavalcante/temperatura-cep/internal/entity"
	"github.com/GeovaneCavalcante/temperatura-cep/internal/infra/web/presenter"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/logger"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/temperature"
)

type WebTemperatureHandler struct {
	FindTemperature entity.FindTemperatureUseCase
}

func NewWebTemperatureHandler(findTemperature entity.FindTemperatureUseCase) *WebTemperatureHandler {
	return &WebTemperatureHandler{
		FindTemperature: findTemperature,
	}
}

func (h *WebTemperatureHandler) TemperatureByCepHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	logger.Info("[TemperatureByCepHandler] starting handler")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	zipcode := r.URL.Query().Get("cep")

	tempData, err := h.FindTemperature.Execute(ctx, zipcode)

	if err == entity.ErrZipCodeNotFound {
		logger.Error("[TemperatureByCepHandler] zipcode not found", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not found zipcode"))
		return
	}

	if err != nil {
		logger.Error("[TemperatureByCepHandler] fail to execute usecase FindTemperature", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	tempDataPresenter := toTemperaturePresenter(tempData)

	json.NewEncoder(w).Encode(&tempDataPresenter)
}

func toTemperaturePresenter(temp *temperature.Info) *presenter.TemperaturePresenter {
	return &presenter.TemperaturePresenter{
		TempC: temp.Celsius,
		TempF: temp.Fahrenheit,
		TempK: temp.Kelvin,
	}
}
