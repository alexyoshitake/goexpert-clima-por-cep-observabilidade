package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"clima-por-cep-observabilidade/internal/providers"
)

func TestServiceBReturnsClimate(t *testing.T) {
	cepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws/01001000/json/" {
			t.Errorf("path = %s, want /ws/01001000/json/", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"localidade":"São Paulo"}`))
	}))
	defer cepServer.Close()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/current.json" {
			t.Errorf("path = %s, want /current.json", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer weatherServer.Close()

	handler := NewHandler(
		providers.NewViaCEP(cepServer.Client(), cepServer.URL),
		providers.NewWeatherAPI(weatherServer.Client(), weatherServer.URL, "test-key"),
	)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/clima-por-cep/01001000", nil)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if got := recorder.Body.String(); got != "{\"city\":\"São Paulo\",\"temp_C\":28.5,\"temp_F\":83.3,\"temp_K\":301.5}\n" {
		t.Fatalf("body = %q, want climate response", got)
	}
}

func TestServiceBReturnsInvalidCEP(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/clima-por-cep/1234567", nil)

	NewHandler(nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", recorder.Code)
	}
	if recorder.Body.String() != "invalid zipcode\n" {
		t.Fatalf("body = %q, want invalid zipcode", recorder.Body.String())
	}
}

func TestServiceBReturnsNotFound(t *testing.T) {
	cepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"erro":"true"}`))
	}))
	defer cepServer.Close()

	handler := NewHandler(
		providers.NewViaCEP(cepServer.Client(), cepServer.URL),
		nil,
	)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/clima-por-cep/99999999", nil)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
	if recorder.Body.String() != "can not find zipcode\n" {
		t.Fatalf("body = %q, want can not find zipcode", recorder.Body.String())
	}
}

func TestServiceBReturnsNotFoundForUnknownRoute(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	NewHandler(nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
}
