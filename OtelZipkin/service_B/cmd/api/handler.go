package api

import (
	"OtelZipkin/service_B/cmd/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

func NewWeatherHandler(service service.WeatherServiceInterface) *WeatherHandler {
	return &WeatherHandler{
		service: service,
	}
}

type WeatherHandler struct {
	service service.WeatherServiceInterface
}

func (h *WeatherHandler) GetWeatherByCityHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cepParam := r.URL.Query().Get("cep")
	if len(cepParam) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	_, err := strconv.Atoi(cepParam)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	weather, err := h.service.GetWeather(ctx, cepParam)
	if err != nil {
		if err == service.ErrCepInvalid {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}
		if err == service.ErrCepNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}
