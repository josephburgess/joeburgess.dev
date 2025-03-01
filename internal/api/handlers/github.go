package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/josephburgess/joeburgess-dev/internal/templates"
)

type GithubHandler struct {
	dataUpdater *templates.DataUpdater
}

func NewGithubHandler(dataUpdater *templates.DataUpdater) *GithubHandler {
	return &GithubHandler{
		dataUpdater: dataUpdater,
	}
}

func (h *GithubHandler) HandleGithubData(w http.ResponseWriter, r *http.Request) {
	data := h.dataUpdater.GetData()

	apiData := map[string]any{
		"repos":      data.GithubRepos,
		"activities": data.GitHubActivities,
		"updated":    data.LastUpdated,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apiData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
