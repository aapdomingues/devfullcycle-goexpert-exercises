package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWheatherByCity(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			http.Error(w, "q is empty", http.StatusBadRequest)
			return
		}

		mockResponse := WheaterApiResponse{
			Location: Location{Name: "Test City"},
			Current:  Current{TempC: 25.0, TempF: 77.0},
		}
		jsonResponse, _ := json.Marshal(mockResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})
	mockServer := httptest.NewServer(mux)
	defer mockServer.Close()

	originalURL := wheaterApiUrl
	wheaterApiUrl = mockServer.URL
	defer func() { wheaterApiUrl = originalURL }()

	testCity := City{Localidade: "Test City"}
	weather, err := getWheatherByCity(context.Background(), testCity)
	if err != nil {
		t.Fatalf("getWheatherByCity returned an error: %v", err)
	}

	if weather.Location.Name != "Test City" {
		t.Errorf("Expected city name 'Test City', got '%s'", weather.Location.Name)
	}
	if weather.Current.TempC != 25.0 {
		t.Errorf("Expected temperature in Celsius 25.0, got %f", weather.Current.TempC)
	}
	if weather.Current.TempF != 77.0 {
		t.Errorf("Expected temperature in Fahrenheit 77.0, got %f", weather.Current.TempF)
	}
}

func TestGetCityByCep(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/12345678/json", func(w http.ResponseWriter, r *http.Request) {
		mockResponse := City{
			Cep:        "12345-678",
			Logradouro: "Test Street",
			Localidade: "Test City",
		}
		jsonResponse, _ := json.Marshal(mockResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})
	mockServer := httptest.NewServer(mux)
	defer mockServer.Close()

	originalURL := viaCepApiUrl
	viaCepApiUrl = mockServer.URL
	defer func() { viaCepApiUrl = originalURL }()

	city, err := getCityByCep(context.Background(), "12345678")
	if err != nil {
		t.Fatalf("getCityByCep returned an error: %v", err)
	}

	if city.Localidade != "Test City" {
		t.Errorf("Expected city 'Test City', got '%s'", city.Localidade)
	}
}

func TestGetWheaterByCityHandler(t *testing.T) {
	viaCepMux := http.NewServeMux()
	viaCepMux.HandleFunc("/12345678/json", func(w http.ResponseWriter, r *http.Request) {
		mockResponse := City{Localidade: "Test City"}
		jsonResponse, _ := json.Marshal(mockResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})
	mockViaCep := httptest.NewServer(viaCepMux)
	defer mockViaCep.Close()

	weatherMux := http.NewServeMux()
	weatherMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mockResponse := WheaterApiResponse{
			Current: Current{TempC: 20, TempF: 68},
		}
		jsonResponse, _ := json.Marshal(mockResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})
	mockWeatherAPI := httptest.NewServer(weatherMux)
	defer mockWeatherAPI.Close()

	originalViaCepApiUrl := viaCepApiUrl
	viaCepApiUrl = mockViaCep.URL
	defer func() { viaCepApiUrl = originalViaCepApiUrl }()

	originalWeatherApiUrl := wheaterApiUrl
	wheaterApiUrl = mockWeatherAPI.URL
	defer func() { wheaterApiUrl = originalWeatherApiUrl }()

	req := httptest.NewRequest("GET", "/?cep=12345678", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetWheaterByCityHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"temp_C":20,"temp_F":68,"temp_K":293}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetWheaterByCityHandler_NoCep(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetWheaterByCityHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
