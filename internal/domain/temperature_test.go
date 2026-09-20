package domain

import "testing"

func TestIsValidCEP(t *testing.T) {
	tests := map[string]bool{
		"01001000":  true,
		"29902555":  true,
		"1234567":   false,
		"123456789": false,
		"1234abcd":  false,
		"":          false,
	}

	for cep, want := range tests {
		if got := IsValidCEP(cep); got != want {
			t.Errorf("IsValidCEP(%q) = %v, want %v", cep, got, want)
		}
	}
}

func TestTemperatureFromCelsius(t *testing.T) {
	got := TemperatureFromCelsius("São Paulo", 28.5)
	want := Temperature{
		City:       "São Paulo",
		Celsius:    28.5,
		Fahrenheit: 83.3,
		Kelvin:     301.5,
	}

	if got != want {
		t.Fatalf("TemperatureFromCelsius() = %#v, want %#v", got, want)
	}
}
