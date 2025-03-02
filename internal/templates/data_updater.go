package templates

import (
	"sync"
	"time"

	"github.com/josephburgess/joeburgess.dev/internal/services/github"
	"github.com/josephburgess/joeburgess.dev/internal/services/weather"
)

type DataUpdater struct {
	data            *PageData
	githubService   *github.Client
	weatherService  *weather.Client
	weatherLocation string
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
	}
}

func (du *DataUpdater) GetData() PageData {
	du.data.mu.RLock()
	defer du.data.mu.RUnlock()

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

func (du *DataUpdater) Update() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		repos, err := du.githubService.FetchRepositories()
		if err != nil {
			return
		}

		du.data.mu.Lock()
		du.data.GithubRepos = repos
		du.data.mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		activities, err := du.githubService.FetchActivity()
		if err != nil {
			return
		}

		du.data.mu.Lock()
		du.data.GitHubActivities = activities
		du.data.mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		if du.weatherLocation != "" {
			weather, err := du.weatherService.FetchWeather(du.weatherLocation)
			if err != nil {
				return
			}

			du.data.mu.Lock()
			du.data.Weather = weather
			du.data.mu.Unlock()
		}
	}()

	wg.Wait()

	now := time.Now()
	du.data.mu.Lock()
	du.data.LastUpdated = now.Format("Jan 02 2006 15:04:05")
	du.data.mu.Unlock()
}

func (du *DataUpdater) StartBackgroundUpdater(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		du.Update()
	}
}
