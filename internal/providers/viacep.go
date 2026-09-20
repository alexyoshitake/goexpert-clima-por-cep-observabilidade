package providers

import (
	"encoding/json"
	"errors"
	"net/http"

	"clima-por-cep-observabilidade/internal/domain"
)

type ViaCEP struct {
	client  *http.Client
	baseURL string
}

var erroViaCEP = errors.New("erro no ViaCEP")

func NewViaCEP(client *http.Client, baseURL string) *ViaCEP {
	return &ViaCEP{
		client:  client,
		baseURL: baseURL,
	}
}

type viaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

func (p *ViaCEP) BuscaCEP(cep string) (string, error) {
	endpoint := p.baseURL + "/ws/" + cep + "/json/"
	response, err := p.client.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", erroViaCEP
	}

	var data viaCEPResponse
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return "", err
	}
	if data.Erro == "true" {
		return "", domain.ErrCEPNaoEncontrado
	}

	return data.Localidade, nil
}
