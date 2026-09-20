package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherAPIBuscaTemperatura(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/current.json" {
			t.Fatalf("path = %s, want /current.json", r.URL.Path)
		}
		if r.URL.Query().Get("key") != "test-key" {
			t.Fatalf("key = %q, want test-key", r.URL.Query().Get("key"))
		}
		if r.URL.Query().Get("q") != "São Paulo" {
			t.Fatalf("q = %q, want São Paulo", r.URL.Query().Get("q"))
		}
		_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	provider := NewWeatherAPI(server.Client(), server.URL, "test-key")
	got, err := provider.BuscaTemperatura("São Paulo")
	if err != nil {
		t.Fatalf("BuscaTemperatura() error = %v", err)
	}
	if got != 28.5 {
		t.Fatalf("BuscaTemperatura() = %v, want 28.5", got)
	}
}
