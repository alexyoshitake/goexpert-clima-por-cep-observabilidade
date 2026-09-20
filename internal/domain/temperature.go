package domain

import "math"

type CEPInput struct {
	CEP string `json:"cep"`
}

type Temperature struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func TemperatureFromCelsius(city string, celsius float64) Temperature {
	return Temperature{
		City:       city,
		Celsius:    celsius,
		Fahrenheit: round(celsius*1.8 + 32),
		Kelvin:     round(celsius + 273),
	}
}

func IsValidCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}

	for i := 0; i < len(cep); i++ {
		if cep[i] < '0' || cep[i] > '9' {
			return false
		}
	}

	return true
}

func round(value float64) float64 {
	return math.Round(value*100) / 100
}
