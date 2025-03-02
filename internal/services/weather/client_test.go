package weather

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/josephburgess/joeburgess-dev/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	apiKey := "test_api_key"
	client := NewClient(apiKey)

	assert.Equal(t, apiKey, client.apiKey)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
}

func TestFetchWeather(t *testing.T) {
	tests := []struct {
		name           string
		apiKey         string
		location       string
		apiURL         string
		setEnvVar      bool
		responseStatus int
		responseBody   string
		expectedError  bool
		validateResult func(*testing.T, *models.WeatherData)
	}{
		{
			name:          "empty API key returns nil",
			apiKey:        "",
			location:      "London",
			expectedError: false,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.Nil(t, data)
			},
		},
		{
			name:           "successful response",
			apiKey:         "test_api_key",
			location:       "London",
			apiURL:         "http://localhost:8080",
			responseStatus: http.StatusOK,
			responseBody: `{
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
			}`,
			expectedError: false,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.NotNil(t, data)
				assert.Equal(t, "London", data.Location)
				assert.Equal(t, 15.5, data.Temperature)
				assert.Equal(t, "Clear", data.Condition)
				assert.Equal(t, "https://openweathermap.org/img/wn/01d@2x.png", data.Icon)
				assert.NotZero(t, data.LastUpdated)
			},
		},
		{
			name:           "custom API URL from env var",
			apiKey:         "test_api_key",
			location:       "Paris",
			apiURL:         "https://custom-api.example.com",
			setEnvVar:      true,
			responseStatus: http.StatusOK,
			responseBody: `{
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
			}`,
			expectedError: false,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.NotNil(t, data)
				assert.Equal(t, "Paris", data.Location)
				assert.Equal(t, 20.0, data.Temperature)
				assert.Equal(t, "Clouds", data.Condition)
				assert.Equal(t, "https://openweathermap.org/img/wn/02d@2x.png", data.Icon)
			},
		},
		// {
		// 	name:           "API returns error status",
		// 	apiKey:         "test_api_key",
		// 	location:       "Invalid Location",
		// 	apiURL:         "http://localhost:8080",
		// 	responseStatus: http.StatusBadRequest,
		// 	responseBody:   `{"error":"Invalid location"}`,
		// 	expectedError:  true,
		// 	validateResult: func(t *testing.T, data *models.WeatherData) {
		// 		assert.Nil(t, data)
		// 	},
		// },
		{
			name:           "invalid JSON response",
			apiKey:         "test_api_key",
			location:       "Berlin",
			apiURL:         "http://localhost:8080",
			responseStatus: http.StatusOK,
			responseBody:   `{invalid json}`,
			expectedError:  true,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.Nil(t, data)
			},
		},
		{
			name:           "missing weather field",
			apiKey:         "test_api_key",
			location:       "Berlin",
			apiURL:         "http://localhost:8080",
			responseStatus: http.StatusOK,
			responseBody:   `{"data": "no weather here"}`,
			expectedError:  true,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.Nil(t, data)
			},
		},
		{
			name:           "missing current field",
			apiKey:         "test_api_key",
			location:       "Berlin",
			apiURL:         "http://localhost:8080",
			responseStatus: http.StatusOK,
			responseBody:   `{"weather": {"forecast": {}}}`,
			expectedError:  true,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.Nil(t, data)
			},
		},
		{
			name:           "missing weather array",
			apiKey:         "test_api_key",
			location:       "Berlin",
			apiURL:         "http://localhost:8080",
			responseStatus: http.StatusOK,
			responseBody:   `{"weather": {"current": {"temp": 10.0}}}`,
			expectedError:  true,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.Nil(t, data)
			},
		},
		{
			name:           "empty weather array",
			apiKey:         "test_api_key",
			location:       "Berlin",
			apiURL:         "http://localhost:8080",
			responseStatus: http.StatusOK,
			responseBody:   `{"weather": {"current": {"temp": 10.0, "weather": []}}}`,
			expectedError:  true,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.Nil(t, data)
			},
		},
		{
			name:           "invalid temp format",
			apiKey:         "test_api_key",
			location:       "Berlin",
			apiURL:         "http://localhost:8080",
			responseStatus: http.StatusOK,
			responseBody:   `{"weather": {"current": {"temp": "hot", "weather": [{"main": "Clear", "icon": "01d"}]}}}`,
			expectedError:  true,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.Nil(t, data)
			},
		},
		{
			name:           "missing condition uses default",
			apiKey:         "test_api_key",
			location:       "Berlin",
			apiURL:         "http://localhost:8080",
			responseStatus: http.StatusOK,
			responseBody:   `{"weather": {"current": {"temp": 10.0, "weather": [{"icon": "01d"}]}}}`,
			expectedError:  false,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.NotNil(t, data)
				assert.Equal(t, "Unknown", data.Condition)
			},
		},
		{
			name:           "missing icon uses default",
			apiKey:         "test_api_key",
			location:       "Berlin",
			apiURL:         "http://localhost:8080",
			responseStatus: http.StatusOK,
			responseBody:   `{"weather": {"current": {"temp": 10.0, "weather": [{"main": "Clear"}]}}}`,
			expectedError:  false,
			validateResult: func(t *testing.T, data *models.WeatherData) {
				assert.NotNil(t, data)
				assert.Equal(t, "https://openweathermap.org/img/wn/01d@2x.png", data.Icon)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup httpmock
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Save original env and restore after test
			origEnv := os.Getenv("BREEZE_API_URL")
			defer os.Setenv("BREEZE_API_URL", origEnv)

			if tc.setEnvVar {
				os.Setenv("BREEZE_API_URL", tc.apiURL)
			} else {
				os.Unsetenv("BREEZE_API_URL")
			}

			if tc.apiKey != "" {
				baseURL := tc.apiURL
				if baseURL == "" {
					baseURL = "http://localhost:8080"
				}

				// Use URL encoding for the location
				encodedLocation := tc.location
				if tc.location == "Invalid Location" {
					encodedLocation = "Invalid%20Location"
				}

				requestURL := baseURL + "/api/weather/" + encodedLocation + "?api_key=" + tc.apiKey + "&units=metric"
				httpmock.RegisterResponder("GET", requestURL,
					httpmock.NewStringResponder(tc.responseStatus, tc.responseBody))
			}

			// Create client and call the method
			client := NewClient(tc.apiKey)
			result, err := client.FetchWeather(tc.location)

			// Validate
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tc.validateResult(t, result)

			// Verify API was called with the right parameters
			if tc.apiKey != "" {
				callCount := httpmock.GetCallCountInfo()
				baseURL := tc.apiURL
				if baseURL == "" {
					baseURL = "http://localhost:8080"
				}

				// Handle URL encoding for verification too
				encodedLocation := tc.location
				if tc.location == "Invalid Location" {
					encodedLocation = "Invalid%20Location"
				}

				requestURL := "GET " + baseURL + "/api/weather/" + encodedLocation + "?api_key=" + tc.apiKey + "&units=metric"
				assert.Equal(t, 1, callCount[requestURL], "Expected one call to the weather API endpoint")
			}
		})
	}
}
