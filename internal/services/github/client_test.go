package github

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	username := "testuser"
	client := NewClient(username)

	assert.Equal(t, username, client.username)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
}

func TestFetchRepositories(t *testing.T) {
	originalExcludes := ReposToExclude
	ReposToExclude = []string{"homebrew-formulae"}
	defer func() { ReposToExclude = originalExcludes }()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	username := "testuser"
	url := "https://api.github.com/users/testuser/repos?sort=updated&per_page=10"

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(http.StatusOK, `[
			{"name": "homebrew-formulae", "description": "Homebrew Formulae", "updated_at": "2023-01-03T00:00:00Z"},
			{"name": "repo1", "description": "Test Repo 1", "updated_at": "2023-01-01T00:00:00Z"},
			{"name": "repo2", "description": "Test Repo 2", "updated_at": "2023-01-02T00:00:00Z"}
		]`))

	client := NewClient(username)
	repos, err := client.FetchRepositories()

	callCount := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, callCount["GET "+url])

	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 2)

	for _, repo := range repos {
		assert.NotEqual(t, "homebrew-formulae", repo.Name)
	}

	assert.Equal(t, "repo2", repos[0].Name)
	assert.Equal(t, "repo1", repos[1].Name)
}

func TestFetchRepositoriesError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	username := "testuser"
	url := "https://api.github.com/users/testuser/repos?sort=updated&per_page=10"

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(http.StatusUnauthorized, `{"message": "Bad credentials"}`))

	client := NewClient(username)
	repos, err := client.FetchRepositories()

	callCount := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, callCount["GET "+url])

	assert.Error(t, err)
	assert.Nil(t, repos)
}

func TestFetchActivity(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	username := "testuser"
	url := "https://api.github.com/users/testuser/events?per_page=10"

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(http.StatusOK, `[
			{"type": "PushEvent", "repo": {"name": "testuser/repo1"}, "created_at": "2023-01-01T00:00:00Z"},
			{"type": "WatchEvent", "repo": {"name": "testuser/repo2"}, "created_at": "2023-01-02T00:00:00Z"}
		]`))

	client := NewClient(username)
	activities, err := client.FetchActivity()

	callCount := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, callCount["GET "+url])

	assert.NoError(t, err)
	assert.NotNil(t, activities)
	assert.Len(t, activities, 2)
	assert.Equal(t, "PushEvent", activities[0].Type)
	assert.Equal(t, "pushed commits to", activities[0].Action)
	assert.Equal(t, "WatchEvent", activities[1].Type)
	assert.Equal(t, "starred", activities[1].Action)
}

func TestFetchActivityError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	username := "josephburgess"
	url := "https://api.github.com/users/josephburgess/events?per_page=10"

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(http.StatusUnauthorized, `{"message": "Bad credentials"}`))

	client := NewClient(username)
	activities, err := client.FetchActivity()

	callCount := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, callCount["GET "+url])

	assert.Error(t, err)
	assert.Nil(t, activities)
}
