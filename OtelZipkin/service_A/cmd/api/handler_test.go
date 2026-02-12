package api

import (
	"OtelZipkin/service_A/cmd/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetWeather(ctx context.Context, cep string) (service.Weather, error) {
	args := m.Called(ctx, cep)
	return args.Get(0).(service.Weather), args.Error(1)
}

func TestWeatherHandler_GetWeatherByCityHandler(t *testing.T) {
	t.Run("should return weather", func(t *testing.T) {
		mockService := new(MockWeatherService)
		handler := NewWeatherHandler(mockService)

		weather := service.Weather{
			City:       "Test City",
			Celsius:    25,
			Fahrenheit: 77,
			Kelvin:     298,
		}

		mockService.On("GetWeather", mock.Anything, "12345678").Return(weather, nil)

		Input := RequestBody{
			Cep: "12345678",
		}

		inputJson, _ := json.Marshal(Input)
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(inputJson))
		rr := httptest.NewRecorder()

		handler.GetWeatherByCityHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var actualResponse service.Weather
		err := json.Unmarshal(rr.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, weather, actualResponse)
	})

	t.Run("should return error when cep is invalid", func(t *testing.T) {
		mockService := new(MockWeatherService)
		handler := NewWeatherHandler(mockService)

		mockService.On("GetWeather", mock.Anything, "12345678").Return(service.Weather{}, service.ErrCepInvalid)

		Input := RequestBody{
			Cep: "12345678",
		}

		inputJson, _ := json.Marshal(Input)
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(inputJson))
		rr := httptest.NewRecorder()

		handler.GetWeatherByCityHandler(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.Equal(t, "invalid zipcode", rr.Body.String())
	})

	t.Run("should return error when cep is not found", func(t *testing.T) {
		mockService := new(MockWeatherService)
		handler := NewWeatherHandler(mockService)

		mockService.On("GetWeather", mock.Anything, "12345678").Return(service.Weather{}, service.ErrCepNotFound)

		Input := RequestBody{
			Cep: "12345678",
		}

		inputJson, _ := json.Marshal(Input)
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(inputJson))
		rr := httptest.NewRecorder()

		handler.GetWeatherByCityHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "can not find zipcode", rr.Body.String())
	})

	t.Run("should return internal server error", func(t *testing.T) {
		mockService := new(MockWeatherService)
		handler := NewWeatherHandler(mockService)

		mockService.On("GetWeather", mock.Anything, "12345678").Return(service.Weather{}, errors.New("internal error"))

		Input := RequestBody{
			Cep: "12345678",
		}

		inputJson, _ := json.Marshal(Input)
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(inputJson))
		rr := httptest.NewRecorder()

		handler.GetWeatherByCityHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should return unprocessable entity when cep is not 8 digits", func(t *testing.T) {
		handler := NewWeatherHandler(nil)

		Input := RequestBody{
			Cep: "12345",
		}

		inputJson, _ := json.Marshal(Input)
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(inputJson))
		rr := httptest.NewRecorder()

		handler.GetWeatherByCityHandler(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.Equal(t, "invalid zipcode", rr.Body.String())
	})

	t.Run("should return error when cep is invalid is not a number", func(t *testing.T) {
		handler := NewWeatherHandler(nil)

		Input := RequestBody{
			Cep: "cep",
		}

		inputJson, _ := json.Marshal(Input)
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(inputJson))
		rr := httptest.NewRecorder()

		handler.GetWeatherByCityHandler(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.Equal(t, "invalid zipcode", rr.Body.String())
	})
}
