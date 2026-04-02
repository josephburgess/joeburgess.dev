package templates

import (
	"sync"
	"time"

	"github.com/josephburgess/joeburgess.dev/internal/models"
	"github.com/josephburgess/joeburgess.dev/internal/services/github"
	"github.com/josephburgess/joeburgess.dev/internal/services/weather"
)

type DataUpdater struct {
	mu              sync.RWMutex
	data            *PageData
	githubService   *github.Client
	weatherService  *weather.Client
	weatherLocation string
	lastUpdated     time.Time
	maxAge          time.Duration
	updating        sync.Mutex
}

func NewDataUpdater(
	githubService *github.Client,
	weatherService *weather.Client,
	weatherLocation string,
	profileImage string,
	githubURL string,
	linkedInURL string,
	breezeURL string,
	email string,
) *DataUpdater {
	data := &PageData{
		ProfileImage: profileImage,
		GithubURL:    githubURL,
		LinkedInURL:  linkedInURL,
		BreezeURL:    breezeURL,
		Email:        email,
		IsDarkMode:   false,
	}

	return &DataUpdater{
		data:            data,
		githubService:   githubService,
		weatherService:  weatherService,
		weatherLocation: weatherLocation,
		maxAge:          1 * time.Hour,
	}
}

func (du *DataUpdater) GetData() PageData {
	if time.Since(du.lastUpdated) > du.maxAge {
		go du.UpdateIfStale()
	}

	du.mu.RLock()
	defer du.mu.RUnlock()

	dataCopy := PageData{
		ProfileImage:     du.data.ProfileImage,
		GithubURL:        du.data.GithubURL,
		LinkedInURL:      du.data.LinkedInURL,
		BreezeURL:        du.data.BreezeURL,
		Email:            du.data.Email,
		IsDarkMode:       du.data.IsDarkMode,
		LastUpdated:      du.data.LastUpdated,
		GithubRepos:      du.data.GithubRepos,
		GitHubActivities: du.data.GitHubActivities,
	}

	if du.data.Weather != nil {
		weatherCopy := *du.data.Weather
		dataCopy.Weather = &weatherCopy
	}

	return dataCopy
}

// UpdateIfStale triggers an update only if one isn't already running.
func (du *DataUpdater) UpdateIfStale() {
	if !du.updating.TryLock() {
		return // another update is already in progress
	}
	defer du.updating.Unlock()

	// double-check after acquiring lock
	if time.Since(du.lastUpdated) <= du.maxAge {
		return
	}

	du.Update()
}

func (du *DataUpdater) Update() {
	var wg sync.WaitGroup

	var repos []models.Repository
	var activities []models.Activity
	var weatherData *models.WeatherData

	wg.Add(3)

	go func() {
		defer wg.Done()
		r, err := du.githubService.FetchRepositories()
		if err != nil {
			return
		}
		repos = r
	}()

	go func() {
		defer wg.Done()
		a, err := du.githubService.FetchActivity()
		if err != nil {
			return
		}
		activities = a
	}()

	go func() {
		defer wg.Done()
		if du.weatherLocation != "" {
			w, err := du.weatherService.FetchWeather(du.weatherLocation)
			if err != nil {
				return
			}
			weatherData = w
		}
	}()

	wg.Wait()

	du.mu.Lock()
	defer du.mu.Unlock()

	if repos != nil {
		du.data.GithubRepos = repos
	}
	if activities != nil {
		du.data.GitHubActivities = activities
	}
	if weatherData != nil {
		du.data.Weather = weatherData
	}

	now := time.Now()
	du.data.LastUpdated = now.Format("Jan 02 2006 15:04:05")
	du.lastUpdated = now
}
