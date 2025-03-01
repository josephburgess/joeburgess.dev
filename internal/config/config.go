package config

import (
	"os"
)

type Config struct {
	ServerAddress   string
	GithubUsername  string
	WeatherLocation string
	WeatherAPIKey   string
	ProfileImage    string
	GithubURL       string
	LinkedInURL     string
	BreezeURL       string
	Email           string
}

func Load() *Config {
	return &Config{
		ServerAddress:   getEnv("SERVER_ADDRESS", ":8081"),
		GithubUsername:  getEnv("GITHUB_USERNAME", "josephburgess"),
		WeatherLocation: getEnv("WEATHER_LOCATION", "Buenosaires, AR"),
		WeatherAPIKey:   os.Getenv("BREEZE_API_KEY"),
		ProfileImage:    "/static/images/profile.png",
		GithubURL:       "https://github.com/josephburgess",
		LinkedInURL:     "https://linkedin.com/in/josephburgessmba",
		BreezeURL:       "https://github.com/josephburgess/breeze",
		Email:           "joe@joeburgess.dev",
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
