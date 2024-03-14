package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GeovaneCavalcante/temperatura-cep/configs"
	"github.com/GeovaneCavalcante/temperatura-cep/internal/entity"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type WebTemperatureProxyHandler struct {
	ProxyTemperature entity.ProxyTemperatureUseCase
}

func NewWebTemperatureProxyHandler(proxyTemperature entity.ProxyTemperatureUseCase) *WebTemperatureProxyHandler {
	return &WebTemperatureProxyHandler{
		ProxyTemperature: proxyTemperature,
	}
}

func (h *WebTemperatureProxyHandler) TemperatureProxyHandler(w http.ResponseWriter, r *http.Request) {
	configs := configs.GetEnvVars()
	tr := otel.Tracer(configs.ServerAName)

	ctx := r.Context()

	logger.Info("[TemperatureByCepHandler] starting handler")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	zipcode := r.URL.Query().Get("cep")

	if !zipCodeValidator(zipcode) {
		logger.Error("[TemperatureByCepHandler] invalid zipcode", nil)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}
	ctx, span := tr.Start(ctx, "response time when calling service B")
	tempData, err := h.ProxyTemperature.Execute(ctx, zipcode)
	span.SetStatus(codes.Ok, "response time when calling service B")
	span.End()
	if err != nil {
		logger.Error("[TemperatureByCepHandler] fail to execute usecase ProxyTemperature", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(&tempData)

}

func zipCodeValidator(zipCode string) bool {
	return len(zipCode) == 8
}
