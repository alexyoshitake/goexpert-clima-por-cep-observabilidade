package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServiceAForwardsCEPUsingGET(t *testing.T) {
	called := false
	bServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/clima-por-cep/01001000" {
			t.Errorf("path = %s, want /clima-por-cep/01001000", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading body: %v", err)
		}
		if len(body) != 0 {
			t.Errorf("body = %q, want empty", body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"city":"São Paulo"}`))
	}))
	defer bServer.Close()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"cep":"01001000"}`),
	)

	NewServiceAHandler(bServer.URL, bServer.Client()).ServeHTTP(recorder, request)

	if !called {
		t.Fatal("Service B was not called")
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if recorder.Body.String() != `{"city":"São Paulo"}` {
		t.Fatalf("body = %q, want response from Service B", recorder.Body.String())
	}
}

func TestServiceARejectsInvalidCEP(t *testing.T) {
	tests := []string{
		`{"cep":"1234567"}`,
		`{"cep":"123456789"}`,
		`{"cep":"1234abcd"}`,
		`{"cep":12345678}`,
		`{"cep":null}`,
		`{"cep":"01001000"} {}`,
	}

	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))

			NewServiceAHandler("http://service-b", nil).ServeHTTP(recorder, request)

			if recorder.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want 422", recorder.Code)
			}
			if recorder.Body.String() != "invalid zipcode\n" {
				t.Fatalf("body = %q, want invalid zipcode", recorder.Body.String())
			}
		})
	}
}

func TestServiceAProxiesServiceBError(t *testing.T) {
	bServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("can not find zipcode\n"))
	}))
	defer bServer.Close()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"cep":"99999999"}`),
	)

	NewServiceAHandler(bServer.URL, bServer.Client()).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
	if recorder.Body.String() != "can not find zipcode\n" {
		t.Fatalf("body = %q, want can not find zipcode", recorder.Body.String())
	}
}
