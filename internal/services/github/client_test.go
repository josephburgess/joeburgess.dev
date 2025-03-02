package github

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/josephburgess/joeburgess-dev/internal/models"
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
	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectedError  bool
		expectedRepos  int
	}{
		{
			name:           "successful response",
			responseStatus: http.StatusOK,
			responseBody: `[
        {"name": "gust", "updated_at": "2023-01-01T00:00:00Z"},
        {"name": "breeze", "updated_at": "2023-01-02T00:00:00Z"}
      ]`,
			expectedError: false,
			expectedRepos: 2,
		},
		{
			name:           "API err",
			responseStatus: http.StatusUnauthorized,
			responseBody:   `{"message": "noway jose"}`,
			expectedError:  true,
			expectedRepos:  0,
		},
		{
			name:           "more than 6 repos",
			responseStatus: http.StatusOK,
			responseBody: `[
				{"name": "madvillainy", "updated_at": "2023-01-07T00:00:00Z"},
				{"name": "borrowed time", "updated_at": "2023-01-06T00:00:00Z"},
				{"name": "clock ticks faster", "updated_at": "2023-01-05T00:00:00Z"},
				{"name": "the hour they", "updated_at": "2023-01-04T00:00:00Z"},
				{"name": "knock the sick blaster", "updated_at": "2023-01-03T00:00:00Z"},
				{"name": "dick dastardly", "updated_at": "2023-01-02T00:00:00Z"},
				{"name": "mutley with sick laughter", "updated_at": "2023-01-01T00:00:00Z"}
			]`,
			expectedError: false,
			expectedRepos: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			username := "testuser"
			url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=6", username)
			httpmock.RegisterResponder("GET", url,
				httpmock.NewStringResponder(tc.responseStatus, tc.responseBody))

			client := NewClient(username)
			repos, err := client.FetchRepositories()

			callCount := httpmock.GetCallCountInfo()
			assert.Equal(t, 1, callCount["GET "+url])

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, repos)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repos)
				assert.Len(t, repos, tc.expectedRepos)

				if len(repos) > 1 {
					for i := range repos[:len(repos)-1] {
						assert.True(t, repos[i].UpdatedAt.After(repos[i+1].UpdatedAt) || repos[i].UpdatedAt.Equal(repos[i+1].UpdatedAt))
					}
				}
			}
		})
	}
}

func TestFetchActivity(t *testing.T) {
	tests := []struct {
		name             string
		responseStatus   int
		responseBody     string
		expectedError    bool
		expectedActivity int
		checkActivity    func(t *testing.T, activities []models.Activity)
	}{
		{
			name:           "successful response",
			responseStatus: http.StatusOK,
			responseBody: `[
				{"type": "PushEvent", "repo": {"name": "catch a throatful"}, "created_at": "2023-01-01T00:00:00Z"},
				{"type": "WatchEvent", "repo": {"name": "fire vocal"}, "created_at": "2023-01-02T00:00:00Z"}
			]`,
			expectedError:    false,
			expectedActivity: 2,
			checkActivity: func(t *testing.T, activities []models.Activity) {
				assert.Equal(t, "PushEvent", activities[0].Type)
				assert.Equal(t, "pushed commits to", activities[0].Action)
				assert.Equal(t, "WatchEvent", activities[1].Type)
				assert.Equal(t, "starred", activities[1].Action)
			},
		},
		{
			name:             "GitHub API error",
			responseStatus:   http.StatusUnauthorized,
			responseBody:     `{"message": "Unauthorized"}`,
			expectedError:    true,
			expectedActivity: 0,
			checkActivity:    nil,
		},
		{
			name:           "more than 6 activities",
			responseStatus: http.StatusOK,
			responseBody: `[
				{"type": "PushEvent", "repo": {"name": "testuser/repo1"}, "created_at": "2023-01-01T00:00:00Z"},
				{"type": "WatchEvent", "repo": {"name": "testuser/repo2"}, "created_at": "2023-01-02T00:00:00Z"},
				{"type": "ForkEvent", "repo": {"name": "testuser/repo3"}, "created_at": "2023-01-03T00:00:00Z"},
				{"type": "CreateEvent", "repo": {"name": "testuser/repo4"}, "created_at": "2023-01-04T00:00:00Z"},
				{"type": "IssuesEvent", "repo": {"name": "testuser/repo5"}, "created_at": "2023-01-05T00:00:00Z"},
				{"type": "PullRequestEvent", "repo": {"name": "testuser/repo6"}, "created_at": "2023-01-06T00:00:00Z"},
				{"type": "IssueCommentEvent", "repo": {"name": "testuser/repo7"}, "created_at": "2023-01-07T00:00:00Z"}
			]`,
			expectedError:    false,
			expectedActivity: 6,
			checkActivity:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			username := "josephburgess"
			url := fmt.Sprintf("https://api.github.com/users/%s/events?per_page=10", username)
			httpmock.RegisterResponder("GET", url,
				httpmock.NewStringResponder(tc.responseStatus, tc.responseBody))

			client := NewClient(username)
			activities, err := client.FetchActivity()

			callCount := httpmock.GetCallCountInfo()
			assert.Equal(t, 1, callCount["GET "+url])

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, activities)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, activities)
				assert.Len(t, activities, tc.expectedActivity)

				if tc.checkActivity != nil {
					tc.checkActivity(t, activities)
				}

				for _, activity := range activities {
					assert.Equal(t, "https://github.com/"+activity.RepoName, activity.URL)
				}
			}
		})
	}
}

func TestMapActivityAction(t *testing.T) {
	tests := []struct {
		eventType      string
		expectedAction string
	}{
		{"PushEvent", "pushed commits to"},
		{"CreateEvent", "created"},
		{"IssuesEvent", "updated an issue in"},
		{"PullRequestEvent", "worked on a pull request in"},
		{"WatchEvent", "starred"},
		{"ForkEvent", "forked"},
		{"IssueCommentEvent", "commented on an issue in"},
		{"UnknownEvent", "worked on"},
	}

	for _, tc := range tests {
		t.Run(tc.eventType, func(t *testing.T) {
			action := models.MapActivityAction(tc.eventType)
			assert.Equal(t, tc.expectedAction, action)
		})
	}
}
