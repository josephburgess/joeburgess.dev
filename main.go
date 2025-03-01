package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/josephburgess/joeburgess-dev/internal/api"
	"github.com/josephburgess/joeburgess-dev/internal/config"
	"github.com/josephburgess/joeburgess-dev/internal/logging"
	"github.com/josephburgess/joeburgess-dev/internal/services/github"
	"github.com/josephburgess/joeburgess-dev/internal/services/weather"
	"github.com/josephburgess/joeburgess-dev/internal/templates"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it. Using environment variables directly.")
	}

	logger := logging.NewLogger()
	defer logger.Sync()

	cfg := config.Load()
	logging.Info("Configuration loaded")

	githubService := github.NewClient(cfg.GithubUsername)
	weatherService := weather.NewClient(cfg.WeatherAPIKey)

	tmplRenderer := templates.NewRenderer()
	dataUpdater := templates.NewDataUpdater(
		githubService,
		weatherService,
		cfg.WeatherLocation,
		cfg.ProfileImage,
		cfg.GithubURL,
		cfg.LinkedInURL,
		cfg.BreezeURL,
		cfg.Email,
	)

	dataUpdater.Update()

	go dataUpdater.StartBackgroundUpdater(15 * time.Minute)

	r := api.Setup(tmplRenderer, dataUpdater)

	logging.Info("Server starting on %s", cfg.ServerAddress)
	if err := r.ListenAndServe(); err != nil {
		logging.Error("Failed to start server", err)
		os.Exit(1)
	}
}
