package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/josephburgess/joeburgess-dev/internal/logging"
	"github.com/josephburgess/joeburgess-dev/internal/models"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) FetchWeather(location string) (*models.WeatherData, error) {
	if c.apiKey == "" {
		return nil, nil
	}

	baseURL := os.Getenv("BREEZE_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	logging.Info("Using BREEZE_API_URL: %s", baseURL)
	requestURL := fmt.Sprintf("%s/api/weather/%s?api_key=%s&units=%s",
		baseURL,
		url.QueryEscape(location),
		c.apiKey,
		"metric",
	)

	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		logging.Error("HTTP request failed", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("breeze API returned status: %s, body: %s", resp.Status, string(bodyBytes))
		logging.Warn(errMsg)
		return nil, errors.New(errMsg)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	weatherResponse, ok := result["weather"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid API response format")
	}

	current, ok := weatherResponse["current"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid API response format")
	}

	weatherArray, ok := current["weather"].([]any)
	if !ok || len(weatherArray) == 0 {
		return nil, fmt.Errorf("weather data not found in API response")
	}

	weather, ok := weatherArray[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid weather data format")
	}

	temp, ok := current["temp"].(float64)
	if !ok {
		return nil, fmt.Errorf("temperature data not found or invalid")
	}

	condition, ok := weather["main"].(string)
	if !ok {
		condition = "Unknown"
	}

	icon, ok := weather["icon"].(string)
	if !ok {
		icon = "01d"
	}

	return &models.WeatherData{
		Location:    location,
		Temperature: temp,
		Condition:   condition,
		Icon:        fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", icon),
		LastUpdated: time.Now(),
	}, nil
}
