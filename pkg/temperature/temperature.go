package temperature

import "context"

type Info struct {
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
	Kelvin     float64 `json:"kelvin"`
}

type TemperatureFetcher interface {
	GetByCity(ctx context.Context, city string) (*Info, error)
}
