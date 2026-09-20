package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"clima-por-cep-observabilidade/internal/domain"
)

func TestViaCEPBuscaCEP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws/01001000/json/" {
			t.Fatalf("path = %s, want /ws/01001000/json/", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"localidade":"São Paulo"}`))
	}))
	defer server.Close()

	provider := NewViaCEP(server.Client(), server.URL)
	got, err := provider.BuscaCEP("01001000")
	if err != nil {
		t.Fatalf("BuscaCEP() error = %v", err)
	}
	if got != "São Paulo" {
		t.Fatalf("BuscaCEP() = %q, want São Paulo", got)
	}
}

func TestViaCEPBuscaCEPNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"erro":"true"}`))
	}))
	defer server.Close()

	provider := NewViaCEP(server.Client(), server.URL)
	_, err := provider.BuscaCEP("99999999")
	if err != domain.ErrCEPNaoEncontrado {
		t.Fatalf("BuscaCEP() error = %v, want %v", err, domain.ErrCEPNaoEncontrado)
	}
}
