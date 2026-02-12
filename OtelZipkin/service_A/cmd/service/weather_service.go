package service

import (
	"OtelZipkin/platform/httpclient"
	"OtelZipkin/service_A/cmd/config"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

type Weather struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func NewWeatherService(cfg *config.Config, client httpclient.HTTPClient) *WeatherService {
	return &WeatherService{
		cfg:    cfg,
		client: client,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, cep string) (Weather, error) {
	wheater, err := s.getWeatherByCep(ctx, cep)
	if err != nil {
		return Weather{}, err
	}

	return wheater, nil
}

func (s *WeatherService) getWeatherByCep(ctx context.Context, cep string) (Weather, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.cfg.ServiceBUrl+"?cep="+cep, nil)
	if err != nil {
		return Weather{}, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return Weather{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return Weather{}, ErrCepNotFound
		}
		return Weather{}, ErrCepInvalid
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Weather{}, err
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return Weather{}, err
	}

	return weather, nil
}
