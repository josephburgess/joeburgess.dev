package weather

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	apiKey := "test_api_key"
	client := NewClient(apiKey)

	assert.Equal(t, apiKey, client.apiKey)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
}

func TestFetchWeatherEmptyAPIKey(t *testing.T) {
	client := NewClient("")
	weatherData, err := client.FetchWeather("London")

	assert.NoError(t, err)
	assert.Nil(t, weatherData)
}

func TestFetchWeatherSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	origEnv := os.Getenv("BREEZE_API_URL")
	defer os.Setenv("BREEZE_API_URL", origEnv)
	os.Unsetenv("BREEZE_API_URL")

	apiKey := "test_api_key"
	location := "London"
	url := "http://localhost:8080/api/weather/London?api_key=test_api_key&units=metric"

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(http.StatusOK, `{
			"weather": {
				"current": {
					"temp": 15.5,
					"weather": [
						{
							"main": "Clear",
							"icon": "01d"
						}
					]
				}
			}
		}`))

	client := NewClient(apiKey)
	weatherData, err := client.FetchWeather(location)

	callCount := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, callCount["GET "+url])

	assert.NoError(t, err)
	assert.NotNil(t, weatherData)
	assert.Equal(t, "London", weatherData.Location)
	assert.Equal(t, 15.5, weatherData.Temperature)
	assert.Equal(t, "Clear", weatherData.Condition)
	assert.Equal(t, "https://openweathermap.org/img/wn/01d@2x.png", weatherData.Icon)
	assert.NotZero(t, weatherData.LastUpdated)
}

func TestFetchWeatherAPIError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	origEnv := os.Getenv("BREEZE_API_URL")
	defer os.Setenv("BREEZE_API_URL", origEnv)
	os.Unsetenv("BREEZE_API_URL") // Use default URL for test

	apiKey := "test_api_key"
	location := "InvalidLocation"
	url := "http://localhost:8080/api/weather/InvalidLocation?api_key=test_api_key&units=metric"

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(http.StatusBadRequest, `{"error": "Location not found"}`))

	client := NewClient(apiKey)
	weatherData, err := client.FetchWeather(location)

	callCount := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, callCount["GET "+url])

	assert.Error(t, err)
	assert.Nil(t, weatherData)
}

func TestFetchWeatherInvalidResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	origEnv := os.Getenv("BREEZE_API_URL")
	defer os.Setenv("BREEZE_API_URL", origEnv)
	os.Unsetenv("BREEZE_API_URL")

	apiKey := "madvillainy"
	location := "London"
	url := "http://localhost:8080/api/weather/London?api_key=madvillainy&units=metric"

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(http.StatusOK, `{"data": "not what we expect"}`))

	client := NewClient(apiKey)
	weatherData, err := client.FetchWeather(location)

	callCount := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, callCount["GET "+url])

	assert.Error(t, err)
	assert.Nil(t, weatherData)
}

func TestFetchWeatherCustomURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	origEnv := os.Getenv("BREEZE_API_URL")
	defer os.Setenv("BREEZE_API_URL", origEnv)

	customURL := "https://custom-api.example.com"
	os.Setenv("BREEZE_API_URL", customURL)

	apiKey := "test_api_key"
	location := "Paris"
	url := "https://custom-api.example.com/api/weather/Paris?api_key=test_api_key&units=metric"

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(http.StatusOK, `{
			"weather": {
				"current": {
					"temp": 20.0,
					"weather": [
						{
							"main": "Clouds",
							"icon": "02d"
						}
					]
				}
			}
		}`))

	client := NewClient(apiKey)
	weatherData, err := client.FetchWeather(location)

	callCount := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, callCount["GET "+url])

	assert.NoError(t, err)
	assert.NotNil(t, weatherData)
	assert.Equal(t, "Paris", weatherData.Location)
	assert.Equal(t, 20.0, weatherData.Temperature)
	assert.Equal(t, "Clouds", weatherData.Condition)
	assert.Equal(t, "https://openweathermap.org/img/wn/02d@2x.png", weatherData.Icon)
}
