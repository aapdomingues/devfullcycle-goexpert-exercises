package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type WheaterApiResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Current struct {
	LastUpdatedEpoch int     `json:"last_updated_epoch"`
	LastUpdated      string  `json:"last_updated"`
	TempC            float64 `json:"temp_c"`
	TempF            float64 `json:"temp_f"`
}

type Wheater struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type City struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

var (
	viaCepApiUrl  = "https://viacep.com.br/ws"
	wheaterApiUrl = "http://api.weatherapi.com/v1/current.json"
	// ?key=9186a23014e94b3d97e01725263101&q=London&aqi=no
	API_KEY = "9186a23014e94b3d97e01725263101" //TODO - Mover para um lugar seguro
)

func GetWheaterByCityHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancelApi := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelApi()

	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	city, err := getCityByCep(ctx, cepParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wheatherApiResponse, err := getWheatherByCity(ctx, city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	wheater := Wheater{
		Celsius:    wheatherApiResponse.Current.TempC,
		Fahrenheit: wheatherApiResponse.Current.TempF,
		Kelvin:     wheatherApiResponse.Current.TempC + 273,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wheater)

}

func getCityByCep(ctx context.Context, cep string) (City, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", viaCepApiUrl+"/"+cep+"/json", nil)
	if err != nil {
		return City{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return City{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return City{}, err
	}

	var city City
	err = json.Unmarshal(body, &city)
	if err != nil {
		return City{}, err
	}

	return city, nil
}

func getWheatherByCity(ctx context.Context, city City) (WheaterApiResponse, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", wheaterApiUrl+"?key="+API_KEY+"&q="+url.QueryEscape(city.Localidade), nil)
	if err != nil {
		return WheaterApiResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return WheaterApiResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return WheaterApiResponse{}, err
	}

	var wheaterApiResponse WheaterApiResponse
	err = json.Unmarshal(body, &wheaterApiResponse)
	if err != nil {
		return WheaterApiResponse{}, err
	}

	return wheaterApiResponse, nil

}
