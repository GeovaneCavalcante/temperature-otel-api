package entity

import (
	"context"
	"errors"

	"github.com/GeovaneCavalcante/temperatura-cep/pkg/temperature"
)

type TemperatureInfo struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

var ErrZipCodeNotFound error = errors.New("zipcode not found")

type FindTemperatureUseCase interface {
	Execute(ctx context.Context, zipCode string) (*temperature.Info, error)
}

type ProxyTemperatureUseCase interface {
	Execute(ctx context.Context, zipCode string) (*TemperatureInfo, error)
}
