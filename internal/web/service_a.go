package web

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"clima-por-cep-observabilidade/internal/domain"
)

type ServiceAHandler struct {
	serviceBURL string
	client      *http.Client
}

func NewServiceAHandler(serviceBURL string, client *http.Client) http.Handler {
	if client == nil {
		client = http.DefaultClient
	}

	return &ServiceAHandler{
		serviceBURL: strings.TrimRight(serviceBURL, "/"),
		client:      client,
	}
}

func (h *ServiceAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var input domain.CEPInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil || !domain.IsValidCEP(input.CEP) {
		http.Error(w, domain.ErrCEPInvalido.Error(), http.StatusUnprocessableEntity)
		return
	}

	var extraPayload any
	if err := decoder.Decode(&extraPayload); err != io.EOF {
		http.Error(w, domain.ErrCEPInvalido.Error(), http.StatusUnprocessableEntity)
		return
	}

	targetURL := h.serviceBURL + "/clima-por-cep/" + url.PathEscape(input.CEP)
	request, err := http.NewRequestWithContext(r.Context(), http.MethodGet, targetURL, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response, err := h.client.Do(request)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if contentType := response.Header.Get("Content-Type"); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(response.StatusCode)
	_, _ = io.Copy(w, response.Body)
}
