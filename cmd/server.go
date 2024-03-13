package main

import (
	"context"
	"log"

	"github.com/GeovaneCavalcante/temperatura-cep/configs"
	"github.com/GeovaneCavalcante/temperatura-cep/internal/infra/web/handler"
	"github.com/GeovaneCavalcante/temperatura-cep/internal/infra/web/webserver"
	"github.com/GeovaneCavalcante/temperatura-cep/internal/usecase/temperature"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/address/viacep"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/opentelemetry"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/requester/resty"
	"github.com/GeovaneCavalcante/temperatura-cep/pkg/temperature/weather"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func newResource(serviceName, version string) (*resource.Resource, error) {
	ctx := context.Background()
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(version),
		))
}

func ServerACommand(config *configs.Conf) *cobra.Command {
	return &cobra.Command{
		Use:   config.ServerAName,
		Short: "Inicialização serviço A",
		Run: func(cmd *cobra.Command, args []string) {
			configs, err := configs.LoadConfig(".")
			if err != nil {
				panic(err)
			}
			otelResource, _ := newResource(configs.ServerAName, configs.Version)

			tp, err := opentelemetry.InitTracer(otelResource)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := tp.Shutdown(context.Background()); err != nil {
					log.Printf("Error shutting down tracer provider: %v", err)
				}
			}()

			requester := resty.New()

			proxyTemperature := temperature.NewProxyTemperatureUseCase(configs.TemperatureUrl, requester)

			ws := webserver.New(configs.WebServerPort)
			tH := handler.NewWebTemperatureProxyHandler(proxyTemperature)
			ws.AddHandler("/temperature", tH.TemperatureProxyHandler)
			ws.Start(":8080")
		},
	}
}

func ServerBCommand(config *configs.Conf) *cobra.Command {
	return &cobra.Command{
		Use:   config.ServerBName,
		Short: "Inicialização serviço B",
		Run: func(cmd *cobra.Command, args []string) {
			configs, err := configs.LoadConfig(".")
			if err != nil {
				panic(err)
			}
			otelResource, _ := newResource(configs.ServerBName, configs.Version)

			tp, err := opentelemetry.InitTracer(otelResource)
			if err != nil {
				log.Fatal(err)
			}
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := tp.Shutdown(context.Background()); err != nil {
					log.Printf("Error shutting down tracer provider: %v", err)
				}
			}()

			requester := resty.New()
			viaCepFetcher := viacep.New(configs.ViaCepApiUrl, requester)
			weatherFetcher := weather.New(configs.WeatherApiUrl, configs.WeatherApiKey, requester)

			findTemperature := temperature.NewFindTemperatureUseCase(viaCepFetcher, weatherFetcher)

			ws := webserver.New(configs.WebServerPort)
			tH := handler.NewWebTemperatureHandler(findTemperature)
			ws.AddHandler("/temperature", tH.TemperatureByCepHandler)
			ws.Start(":8081")
		},
	}
}

func main() {

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(ServerACommand(configs))
	rootCmd.AddCommand(ServerBCommand(configs))
	rootCmd.Execute()
}
