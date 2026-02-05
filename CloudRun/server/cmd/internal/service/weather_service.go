package service

import (
	"CloudRun/server/cmd/internal/config"
	"CloudRun/server/cmd/internal/httpclient"

	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

var (
	ErrCepInvalid         = errors.New("invalid zipcode")
	ErrCepNotFound        = errors.New("can not find zipcode")
	ErrWeatherCityInvalid = errors.New("invalid city")
)

type WeatherServiceInterface interface {
	GetWeather(ctx context.Context, cep string) (Weather, error)
}

type WeatherService struct {
	cfg    *config.Config
	client httpclient.HTTPClient
}

type WeatherApiResponse struct {
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

type Weather struct {
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
	Erro        string `json:"erro"`
}

func NewWeatherService(cfg *config.Config, client httpclient.HTTPClient) *WeatherService {
	return &WeatherService{
		cfg:    cfg,
		client: client,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, cep string) (Weather, error) {
	city, err := s.getCityByCep(ctx, cep)
	if err != nil {
		return Weather{}, err
	}

	wheatherApiResponse, err := s.getWeatherByCity(ctx, city)
	if err != nil {
		return Weather{}, err
	}
	wheater := Weather{
		Celsius:    wheatherApiResponse.Current.TempC,
		Fahrenheit: wheatherApiResponse.Current.TempF,
		Kelvin:     wheatherApiResponse.Current.TempC + 273,
	}

	return wheater, nil
}

func (s *WeatherService) getCityByCep(ctx context.Context, cep string) (City, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.cfg.ViaCepApiUrl+"/"+cep+"/json", nil)
	if err != nil {
		return City{}, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return City{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return City{}, ErrCepInvalid
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return City{}, err
	}

	var city City
	err = json.Unmarshal(body, &city)
	if err != nil {
		return City{}, err
	}

	if city.Erro == "true" {
		return City{}, ErrCepNotFound
	}

	return city, nil
}

func (s *WeatherService) getWeatherByCity(ctx context.Context, city City) (WeatherApiResponse, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", s.cfg.WeatherApiUrl+"?key="+s.cfg.ApiKey+"&q="+url.QueryEscape(city.Localidade), nil)
	if err != nil {
		return WeatherApiResponse{}, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return WeatherApiResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return WeatherApiResponse{}, ErrWeatherCityInvalid
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return WeatherApiResponse{}, err
	}

	var wheaterApiResponse WeatherApiResponse
	err = json.Unmarshal(body, &wheaterApiResponse)
	if err != nil {
		return WeatherApiResponse{}, err
	}

	return wheaterApiResponse, nil

}
