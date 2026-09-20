package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"clima-por-cep-observabilidade/internal/config"
	"clima-por-cep-observabilidade/internal/httpserver"
	"clima-por-cep-observabilidade/internal/providers"
	"clima-por-cep-observabilidade/internal/telemetry"
	"clima-por-cep-observabilidade/internal/web"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	shutdownTelemetry, err := telemetry.InitProvider(
		ctx,
		cfg.OTELServiceName,
		cfg.OTELExporterOTLPEndpoint,
	)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: cfg.HTTPTimeout}
	handler := web.NewHandler(
		providers.NewViaCEP(client, cfg.ViaCEPBaseURL),
		providers.NewWeatherAPI(
			client,
			cfg.WeatherAPIBaseURL,
			cfg.WeatherAPIKey,
		),
	)
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: otelhttp.NewHandler(handler, "service-b"),
	}

	if err := httpserver.Run(ctx, server, shutdownTelemetry); err != nil {
		log.Fatal(err)
	}
}
