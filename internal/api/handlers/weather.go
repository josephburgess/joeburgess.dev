package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/josephburgess/joeburgess.dev/internal/services/weather"
)

type WeatherHandler struct {
	weatherService *weather.Client
}

func NewWeatherHandler(weatherService *weather.Client) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

func (h *WeatherHandler) HandleWeatherData(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "Location parameter is required", http.StatusBadRequest)
		return
	}

	weatherData, err := h.weatherService.FetchWeather(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if weatherData == nil {
		http.Error(w, "Weather API key not configured", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weatherData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
