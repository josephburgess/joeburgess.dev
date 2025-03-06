package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/josephburgess/joeburgess.dev/internal/models"
)

type Client struct {
	username   string
	httpClient *http.Client
}

func NewClient(username string) *Client {
	return &Client{
		username: username,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) FetchRepositories() ([]models.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=10", c.username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var repos []models.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	filteredRepos := make([]models.Repository, 0)
	for _, repo := range repos {
		if ShouldIncludeRepo(repo.Name) {
			filteredRepos = append(filteredRepos, repo)
		}
	}

	sort.Slice(filteredRepos, func(i, j int) bool {
		return filteredRepos[i].UpdatedAt.After(filteredRepos[j].UpdatedAt)
	})

	if len(filteredRepos) > 6 {
		filteredRepos = filteredRepos[:6]
	}

	return filteredRepos, nil
}

func (c *Client) FetchActivity() ([]models.Activity, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events?per_page=10", c.username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var events models.GitHubEventResponse
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	activities := make([]models.Activity, 0)

	for _, event := range events {
		activity := models.Activity{
			Type:      event.Type,
			RepoName:  event.Repo.Name,
			CreatedAt: event.CreatedAt,
			URL:       fmt.Sprintf("https://github.com/%s", event.Repo.Name),
			Action:    models.MapActivityAction(event.Type),
		}

		activities = append(activities, activity)
	}

	if len(activities) > 6 {
		activities = activities[:6]
	}

	return activities, nil
}
