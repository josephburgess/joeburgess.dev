package github

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/josephburgess/joeburgess-dev/internal/models"
	"github.com/stretchr/testify/assert"
)

type MockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

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
				{"name": "repo1", "description": "test repo 1", "html_url": "https://github.com/testuser/repo1", "language": "Go", "stargazers_count": 10, "forks_count": 5, "updated_at": "2023-01-01T00:00:00Z"},
				{"name": "repo2", "description": "test repo 2", "html_url": "https://github.com/testuser/repo2", "language": "JavaScript", "stargazers_count": 20, "forks_count": 10, "updated_at": "2023-01-02T00:00:00Z"}
			]`,
			expectedError: false,
			expectedRepos: 2,
		},
		{
			name:           "GitHub API error",
			responseStatus: http.StatusUnauthorized,
			responseBody:   `{"message": "Unauthorized"}`,
			expectedError:  true,
			expectedRepos:  0,
		},
		{
			name:           "more than 6 repositories",
			responseStatus: http.StatusOK,
			responseBody: `[
				{"name": "repo1", "updated_at": "2023-01-07T00:00:00Z"},
				{"name": "repo2", "updated_at": "2023-01-06T00:00:00Z"},
				{"name": "repo3", "updated_at": "2023-01-05T00:00:00Z"},
				{"name": "repo4", "updated_at": "2023-01-04T00:00:00Z"},
				{"name": "repo5", "updated_at": "2023-01-03T00:00:00Z"},
				{"name": "repo6", "updated_at": "2023-01-02T00:00:00Z"},
				{"name": "repo7", "updated_at": "2023-01-01T00:00:00Z"}
			]`,
			expectedError: false,
			expectedRepos: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &Client{
				username: "testuser",
				httpClient: &http.Client{
					Transport: &MockTransport{
						RoundTripFunc: func(req *http.Request) (*http.Response, error) {
							assert.Contains(t, req.URL.String(), "/users/testuser/repos")
							assert.Equal(t, "GET", req.Method)
							assert.Equal(t, "application/vnd.github.v3+json", req.Header.Get("Accept"))

							return &http.Response{
								StatusCode: tc.responseStatus,
								Body:       io.NopCloser(bytes.NewBufferString(tc.responseBody)),
								Header:     make(http.Header),
							}, nil
						},
					},
				},
			}

			repos, err := client.FetchRepositories()

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, repos)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repos)
				assert.Len(t, repos, tc.expectedRepos)

				if len(repos) > 1 {
					for i := 0; i < len(repos)-1; i++ {
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
				{"type": "PushEvent", "repo": {"name": "testuser/repo1"}, "created_at": "2023-01-01T00:00:00Z"},
				{"type": "WatchEvent", "repo": {"name": "testuser/repo2"}, "created_at": "2023-01-02T00:00:00Z"}
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
			client := &Client{
				username: "testuser",
				httpClient: &http.Client{
					Transport: &MockTransport{
						RoundTripFunc: func(req *http.Request) (*http.Response, error) {
							assert.Contains(t, req.URL.String(), "/users/testuser/events")
							assert.Equal(t, "GET", req.Method)
							assert.Equal(t, "application/vnd.github.v3+json", req.Header.Get("Accept"))

							return &http.Response{
								StatusCode: tc.responseStatus,
								Body:       io.NopCloser(bytes.NewBufferString(tc.responseBody)),
								Header:     make(http.Header),
							}, nil
						},
					},
				},
			}

			activities, err := client.FetchActivity()

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
