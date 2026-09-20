package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"clima-por-cep-observabilidade/internal/domain"
	"clima-por-cep-observabilidade/internal/providers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const routePrefix = "/clima-por-cep/"

func NewHandler(cep *providers.ViaCEP, weather *providers.WeatherAPI) http.HandlerFunc {
	tracer := otel.Tracer("service-b")

	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, routePrefix) {
			http.NotFound(w, r)
			return
		}

		cepValue := strings.TrimPrefix(r.URL.Path, routePrefix)
		if !domain.IsValidCEP(cepValue) {
			http.Error(w, domain.ErrCEPInvalido.Error(), http.StatusUnprocessableEntity)
			return
		}

		_, span := tracer.Start(r.Context(), "Busca CEP")
		city, err := cep.BuscaCEP(cepValue)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
		if errors.Is(err, domain.ErrCEPNaoEncontrado) {
			http.Error(w, domain.ErrCEPNaoEncontrado.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, span = tracer.Start(r.Context(), "Busca temperatura")
		celsius, err := weather.BuscaTemperatura(city)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(domain.TemperatureFromCelsius(city, celsius))
	}
}
