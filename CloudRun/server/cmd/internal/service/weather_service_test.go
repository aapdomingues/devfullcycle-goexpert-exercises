package service

import (
	"CloudRun/server/cmd/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestWeatherService_GetWeather(t *testing.T) {
	cfg := &config.Config{
		ViaCepApiUrl:  "http://test.com",
		WeatherApiUrl: "http://test.com",
		ApiKey:        "test",
	}

	t.Run("should return wheater", func(t *testing.T) {
		mockClient := new(MockHTTPClient)
		service := NewWeatherService(cfg, mockClient)

		city := City{
			Localidade: "Test City",
		}
		cityJson, _ := json.Marshal(city)

		wheater := WeatherApiResponse{
			Current: Current{
				TempC: 25,
			},
		}
		wheaterJson, _ := json.Marshal(wheater)

		mockClient.On("Do", mock.Anything).Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(cityJson)),
			}, nil,
		).Once()

		mockClient.On("Do", mock.Anything).Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(wheaterJson)),
			}, nil,
		).Once()

		wheaterResponse, err := service.GetWeather(context.Background(), "12345678")
		assert.NoError(t, err)
		assert.Equal(t, 25.0, wheaterResponse.Celsius)
		assert.Equal(t, 298.0, wheaterResponse.Kelvin)
	})

	t.Run("should return error when cep is invalid", func(t *testing.T) {
		mockClient := new(MockHTTPClient)
		service := NewWeatherService(cfg, mockClient)

		mockClient.On("Do", mock.Anything).Return(
			&http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil,
		).Once()

		_, err := service.GetWeather(context.Background(), "12345678")
		assert.Error(t, err)
		assert.Equal(t, ErrCepInvalid, err)
	})

	t.Run("should return error when cep is not found", func(t *testing.T) {
		mockClient := new(MockHTTPClient)
		service := NewWeatherService(cfg, mockClient)

		city := City{
			Erro: true,
		}
		cityJson, _ := json.Marshal(city)

		mockClient.On("Do", mock.Anything).Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(cityJson)),
			}, nil,
		).Once()

		_, err := service.GetWeather(context.Background(), "12345678")
		assert.Error(t, err)
		assert.Equal(t, ErrCepNotFound, err)
	})

	t.Run("should return error when city is invalid", func(t *testing.T) {
		mockClient := new(MockHTTPClient)
		service := NewWeatherService(cfg, mockClient)

		city := City{
			Localidade: "Test City",
		}
		cityJson, _ := json.Marshal(city)

		mockClient.On("Do", mock.Anything).Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(cityJson)),
			}, nil,
		).Once()

		mockClient.On("Do", mock.Anything).Return(
			&http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil,
		).Once()

		_, err := service.GetWeather(context.Background(), "12345678")
		assert.Error(t, err)
		assert.Equal(t, ErrWeatherCityInvalid, err)
	})
}
