package providers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type WeatherAPI struct {
	client  *http.Client
	baseURL string
	key     string
}

var erroWeatherAPI = errors.New("erro no WeatherAPI")

func NewWeatherAPI(client *http.Client, baseURL, key string) *WeatherAPI {
	return &WeatherAPI{
		client:  client,
		baseURL: baseURL,
		key:     key,
	}
}

type weatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (p *WeatherAPI) BuscaTemperatura(city string) (float64, error) {
	query := url.Values{}
	query.Set("key", p.key)
	query.Set("q", city)
	endpoint := p.baseURL + "/current.json?" + query.Encode()

	response, err := p.client.Get(endpoint)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, erroWeatherAPI
	}

	var data weatherResponse
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
