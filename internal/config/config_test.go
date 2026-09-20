package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
)

var configEnvironmentKeys = []string{
	"PORT",
	"SERVICE_B_URL",
	"OTEL_SERVICE_NAME",
	"OTEL_EXPORTER_OTLP_ENDPOINT",
	"HTTP_TIMEOUT",
	"VIACEP_BASE_URL",
	"WEATHER_API_BASE_URL",
	"WEATHER_API_KEY",
}

func prepareConfigTest(t *testing.T) string {
	t.Helper()

	viper.Reset()
	t.Cleanup(viper.Reset)

	for _, key := range configEnvironmentKeys {
		t.Setenv(key, "")
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	temporaryDirectory := t.TempDir()
	if err := os.Chdir(temporaryDirectory); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	return temporaryDirectory
}

func writeConfigFile(t *testing.T, directory, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(directory, ".env"), []byte(content), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}
}

func TestLoadReturnsEmptyValuesWhenConfigFileIsMissing(t *testing.T) {
	directory := prepareConfigTest(t)

	got, err := LoadConfig(directory)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := Config{}

	if *got != want {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}
}

func TestLoadReadsDotEnv(t *testing.T) {
	directory := prepareConfigTest(t)
	writeConfigFile(t, directory, `
PORT=9090
SERVICE_B_URL=http://service-b:8081
OTEL_SERVICE_NAME=service-from-file
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
HTTP_TIMEOUT=3s
VIACEP_BASE_URL=https://viacep.example
WEATHER_API_BASE_URL=https://weather.example/v1
WEATHER_API_KEY=file-key
`)

	got, err := LoadConfig(directory)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := Config{
		Port:                     "9090",
		ServiceBURL:              "http://service-b:8081",
		OTELServiceName:          "service-from-file",
		OTELExporterOTLPEndpoint: "otel-collector:4317",
		HTTPTimeout:              3 * time.Second,
		ViaCEPBaseURL:            "https://viacep.example",
		WeatherAPIBaseURL:        "https://weather.example/v1",
		WeatherAPIKey:            "file-key",
	}

	if *got != want {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}
}

func TestLoadEnvironmentOverridesDotEnv(t *testing.T) {
	directory := prepareConfigTest(t)
	writeConfigFile(t, directory, `
PORT=9090
OTEL_SERVICE_NAME=service-from-file
HTTP_TIMEOUT=3s
WEATHER_API_KEY=file-key
`)

	t.Setenv("PORT", "7070")
	t.Setenv("OTEL_SERVICE_NAME", "service-from-environment")
	t.Setenv("HTTP_TIMEOUT", "2s")
	t.Setenv("WEATHER_API_KEY", "environment-key")

	got, err := LoadConfig(directory)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got.Port != "7070" {
		t.Errorf("Port = %q, want 7070", got.Port)
	}
	if got.OTELServiceName != "service-from-environment" {
		t.Errorf("OTELServiceName = %q, want service-from-environment", got.OTELServiceName)
	}
	if got.HTTPTimeout != 2*time.Second {
		t.Errorf("HTTPTimeout = %s, want 2s", got.HTTPTimeout)
	}
	if got.WeatherAPIKey != "environment-key" {
		t.Errorf("WeatherAPIKey = %q, want environment-key", got.WeatherAPIKey)
	}
	if got.ServiceBURL != "" {
		t.Errorf("ServiceBURL = %q, want empty value", got.ServiceBURL)
	}
}
