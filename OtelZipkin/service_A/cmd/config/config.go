package config

import (
	"os"
)

type Config struct {
	ServiceBUrl string
}

func LoadConfig() *Config {
	return &Config{
		ServiceBUrl: getEnv("SERVICE_B_API_URL", "http://localhost:8081/weather"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
