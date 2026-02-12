package config

import (
	"os"
)

type Config struct {
	ViaCepApiUrl  string
	WeatherApiUrl string
	ApiKey        string
}

func LoadConfig() *Config {
	return &Config{
		ViaCepApiUrl:  getEnv("VIACEP_API_URL", "https://viacep.com.br/ws"),
		WeatherApiUrl: getEnv("WHEATER_API_URL", "http://api.weatherapi.com/v1/current.json"),
		ApiKey:        getEnv("API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
