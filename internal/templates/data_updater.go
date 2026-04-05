package templates

import (
	"sync"
	"time"

	"github.com/josephburgess/joeburgess.dev/internal/logging"
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
	return &DataUpdater{
		data: &PageData{
			ProfileImage: profileImage,
			GithubURL:    githubURL,
			LinkedInURL:  linkedInURL,
			BreezeURL:    breezeURL,
			Email:        email,
		},
		githubService:   githubService,
		weatherService:  weatherService,
		weatherLocation: weatherLocation,
		maxAge:          1 * time.Hour,
	}
}

func (du *DataUpdater) GetData() PageData {
	du.mu.RLock()
	stale := time.Since(du.lastUpdated) > du.maxAge
	data := du.copyData()
	du.mu.RUnlock()

	if stale {
		go du.UpdateIfStale()
	}

	return data
}

// UpdateIfStale triggers an update only if one isn't already running.
func (du *DataUpdater) UpdateIfStale() {
	if !du.updating.TryLock() {
		return
	}
	defer du.updating.Unlock()

	du.mu.RLock()
	fresh := time.Since(du.lastUpdated) <= du.maxAge
	du.mu.RUnlock()

	if fresh {
		return
	}

	du.Update()
}

func (du *DataUpdater) Update() {
	var (
		wg           sync.WaitGroup
		repos        []models.Repository
		activities   []models.Activity
		weatherData  *models.WeatherData
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		r, err := du.githubService.FetchRepositories()
		if err != nil {
			logging.Error("Failed to fetch repositories", err)
			return
		}
		repos = r
	}()

	go func() {
		defer wg.Done()
		a, err := du.githubService.FetchActivity()
		if err != nil {
			logging.Error("Failed to fetch GitHub activity", err)
			return
		}
		activities = a
	}()

	go func() {
		defer wg.Done()
		if du.weatherLocation == "" {
			return
		}
		w, err := du.weatherService.FetchWeather(du.weatherLocation)
		if err != nil {
			logging.Error("Failed to fetch weather", err)
			return
		}
		weatherData = w
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

func (du *DataUpdater) copyData() PageData {
	d := PageData{
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
		d.Weather = &weatherCopy
	}
	return d
}
