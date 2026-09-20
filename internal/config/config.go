package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port                     string
	ServiceBURL              string
	OTELServiceName          string
	OTELExporterOTLPEndpoint string
	HTTPTimeout              time.Duration
	ViaCEPBaseURL            string
	WeatherAPIBaseURL        string
	WeatherAPIKey            string
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(arquivoDeConfiguracao(path))
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to read configuration: %w", err)
		}
	}

	return &Config{
		Port:                     viper.GetString("PORT"),
		ServiceBURL:              viper.GetString("SERVICE_B_URL"),
		OTELServiceName:          viper.GetString("OTEL_SERVICE_NAME"),
		OTELExporterOTLPEndpoint: viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"),
		HTTPTimeout:              viper.GetDuration("HTTP_TIMEOUT"),
		ViaCEPBaseURL:            viper.GetString("VIACEP_BASE_URL"),
		WeatherAPIBaseURL:        viper.GetString("WEATHER_API_BASE_URL"),
		WeatherAPIKey:            viper.GetString("WEATHER_API_KEY"),
	}, nil
}

func arquivoDeConfiguracao(path string) string {
	diretorio, err := filepath.Abs(path)
	if err != nil {
		return filepath.Join(path, ".env")
	}

	for {
		arquivo := filepath.Join(diretorio, ".env")
		if _, err := os.Stat(arquivo); err == nil {
			return arquivo
		}

		pai := filepath.Dir(diretorio)
		if pai == diretorio {
			return arquivo
		}
		diretorio = pai
	}
}
