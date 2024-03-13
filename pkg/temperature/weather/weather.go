package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/GeovaneCavalcante/temperatura-cep/pkg/logger"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/requester"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/temperature"
)

type Response struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherFetcher struct {
	Url       string
	Key       string
	Requester requester.Sender
}

func New(url, key string, req requester.Sender) *WeatherFetcher {
	return &WeatherFetcher{
		Url:       url,
		Key:       key,
		Requester: req,
	}
}

func (w *WeatherFetcher) GetByCity(ctx context.Context, city string) (*temperature.Info, error) {

	logger.Info("[WeatherFetcher] starting api for city: " + city)
	cfg := requester.Configuration{
		Url:        fmt.Sprintf("%s?key=%s&q=%s", w.Url, w.Key, url.QueryEscape(city)),
		Method:     "get",
		ContetType: "application/json",
	}

	res, err := w.Requester.Send(ctx, cfg)

	if err != nil {
		logger.Error("[WeatherFetcher] fail to execute http request for city "+city, err)
		return nil, fmt.Errorf("fail to execute http request for city %s: %w", city, err)
	}

	var weatherInfo Response

	err = json.Unmarshal(res.Body, &weatherInfo)
	if err != nil {
		logger.Error("[WeatherFetcher] fail to unmarshal response body for city "+city, err)
		return nil, fmt.Errorf("fail to unmarshal response body for city %s: %w", city, err)
	}

	temperatureInfo := toTemperatureInfo(weatherInfo)

	logger.Info("[WeatherFetcher] finishing api for city " + city)
	return temperatureInfo, nil
}

func toTemperatureInfo(weatherInfo Response) *temperature.Info {
	return &temperature.Info{
		Celsius:    weatherInfo.Current.TempC,
		Kelvin:     celsiusToKelvin(weatherInfo.Current.TempC),
		Fahrenheit: weatherInfo.Current.TempF,
	}
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}
