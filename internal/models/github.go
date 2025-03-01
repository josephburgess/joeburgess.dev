package models

import (
	"encoding/json"
	"time"
)

type Repository struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	URL         string    `json:"html_url"`
	Language    string    `json:"language"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Activity struct {
	Type      string    `json:"type"`
	RepoName  string    `json:"repo"`
	CreatedAt time.Time `json:"created_at"`
	Action    string    `json:"action"`
	URL       string    `json:"url"`
}

type GitHubEventResponse []struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt time.Time `json:"created_at"`
	Payload   json.RawMessage
}

func MapActivityAction(eventType string) string {
	switch eventType {
	case "PushEvent":
		return "pushed commits to"
	case "CreateEvent":
		return "created"
	case "IssuesEvent":
		return "updated an issue in"
	case "PullRequestEvent":
		return "worked on a pull request in"
	case "WatchEvent":
		return "starred"
	case "ForkEvent":
		return "forked"
	case "IssueCommentEvent":
		return "commented on an issue in"
	default:
		return "worked on"
	}
}
