package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clima-por-cep-observabilidade/internal/config"
	"clima-por-cep-observabilidade/internal/httpserver"
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

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   cfg.HTTPTimeout,
	}
	handler := web.NewServiceAHandler(
		cfg.ServiceBURL,
		client,
	)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           otelhttp.NewHandler(handler, "service-a"),
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := httpserver.Run(ctx, server, shutdownTelemetry); err != nil {
		log.Fatal(err)
	}
}
