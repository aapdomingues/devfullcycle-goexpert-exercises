package service

import (
	"OtelZipkin/service_A/cmd/config"
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
		ServiceBUrl: "http://test.com",
	}

	t.Run("should return wheater", func(t *testing.T) {
		mockClient := new(MockHTTPClient)
		service := NewWeatherService(cfg, mockClient)

		wheater := Weather{
			City:       "Test City",
			Celsius:    25.0,
			Kelvin:     298.0,
			Fahrenheit: 77.0,
		}

		wheaterJson, _ := json.Marshal(wheater)

		mockClient.On("Do", mock.Anything).Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(wheaterJson)),
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
		assert.Equal(t, 77.0, wheaterResponse.Fahrenheit)
		assert.Equal(t, "Test City", wheaterResponse.City)
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
}
