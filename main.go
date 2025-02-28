package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type PageData struct {
	ProfileImage     string
	GithubURL        string
	LinkedInURL      string
	Email            string
	IsDarkMode       bool
	GithubRepos      []Repository
	GitHubActivities []Activity
	LastUpdated      string
	Weather          *WeatherData
}

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
	Action    string
	URL       string
}

type WeatherData struct {
	Location    string
	Temperature float64
	Condition   string
	Icon        string
	LastUpdated time.Time
}

type GithubEventResponse []struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt time.Time `json:"created_at"`
	Payload   json.RawMessage
}

func fetchGithubRepos(username string) ([]Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=6", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].UpdatedAt.After(repos[j].UpdatedAt)
	})

	if len(repos) > 6 {
		repos = repos[:6]
	}

	return repos, nil
}

func fetchGithubActivity(username string) ([]Activity, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events?per_page=10", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var events GithubEventResponse
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	activities := make([]Activity, 0)

	for _, event := range events {
		activity := Activity{
			Type:      event.Type,
			RepoName:  event.Repo.Name,
			CreatedAt: event.CreatedAt,
			URL:       fmt.Sprintf("https://github.com/%s", event.Repo.Name),
		}

		switch event.Type {
		case "PushEvent":
			activity.Action = "pushed commits to"
		case "CreateEvent":
			activity.Action = "created"
		case "IssuesEvent":
			activity.Action = "updated an issue in"
		case "PullRequestEvent":
			activity.Action = "worked on a pull request in"
		case "WatchEvent":
			activity.Action = "starred"
		case "ForkEvent":
			activity.Action = "forked"
		case "IssueCommentEvent":
			activity.Action = "commented on an issue in"
		default:
			activity.Action = "worked on"
		}

		activities = append(activities, activity)
	}

	if len(activities) > 6 {
		activities = activities[:6]
	}

	return activities, nil
}

const (
	BuenosAiresLat = -34.6037
	BuenosAiresLon = -58.3816
)

func fetchWeather(location string, apiKey string) (*WeatherData, error) {
	if apiKey == "" {
		return nil, nil
	}

	lat, lon := BuenosAiresLat, BuenosAiresLon

	url := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%.6f&lon=%.6f&appid=%s&units=metric&exclude=minutely,hourly,daily,alerts",
		lat, lon, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %s", resp.Status)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	current, ok := result["current"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid API response format")
	}

	weatherArray, ok := current["weather"].([]any)
	if !ok || len(weatherArray) == 0 {
		return nil, fmt.Errorf("weather data not found in API response")
	}

	weather, ok := weatherArray[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid weather data format")
	}

	temp, ok := current["temp"].(float64)
	if !ok {
		return nil, fmt.Errorf("temperature data not found or invalid")
	}

	condition, ok := weather["main"].(string)
	if !ok {
		condition = "Unknown"
	}

	icon, ok := weather["icon"].(string)
	if !ok {
		icon = "01d"
	}

	return &WeatherData{
		Location:    location,
		Temperature: temp,
		Condition:   condition,
		Icon:        fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", icon),
		LastUpdated: time.Now(),
	}, nil
}

func startBackgroundUpdater(data *PageData, githubUsername string, weatherLocation string, apiKey string, mu *sync.Mutex) {
	go func() {
		for {
			updateData(data, githubUsername, weatherLocation, apiKey, mu)
			time.Sleep(15 * time.Minute)
		}
	}()
}

func updateData(data *PageData, githubUsername string, weatherLocation string, apiKey string, mu *sync.Mutex) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		repos, err := fetchGithubRepos(githubUsername)
		if err != nil {
			log.Printf("Error fetching GitHub repos: %v", err)
			return
		}

		mu.Lock()
		data.GithubRepos = repos
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		activities, err := fetchGithubActivity(githubUsername)
		if err != nil {
			log.Printf("Error fetching GitHub activity: %v", err)
			return
		}

		mu.Lock()
		data.GitHubActivities = activities
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		if weatherLocation != "" && apiKey != "" {
			weather, err := fetchWeather(weatherLocation, apiKey)
			if err != nil {
				log.Printf("Error fetching weather: %v", err)
				return
			}

			mu.Lock()
			data.Weather = weather
			mu.Unlock()
		}
	}()

	wg.Wait()

	now := time.Now()
	mu.Lock()
	data.LastUpdated = now.Format("Jan 02 2006 15:04:05")
	mu.Unlock()

	log.Println("Data updated at:", now.Format(time.RFC3339))
}

func formatDate(t time.Time) string {
	return t.Format("Jan 02, 2006")
}

func timeSince(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff.Hours() < 24 {
		hours := int(diff.Hours())
		if hours < 1 {
			minutes := int(diff.Minutes())
			if minutes < 1 {
				return "just now"
			}
			return fmt.Sprintf("%d minutes ago", minutes)
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff.Hours() < 48 {
		return "yesterday"
	} else if diff.Hours() < 24*7 {
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	} else if diff.Hours() < 24*30 {
		return fmt.Sprintf("%d weeks ago", int(diff.Hours()/(24*7)))
	} else {
		return t.Format("Jan 02")
	}
}

func main() {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"formatDate": formatDate,
		"timeSince":  timeSince,
		"toLower":    strings.ToLower,
	}).ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	err = godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it. Using environment variables directly.")
	}

	weatherApiKey := os.Getenv("OPENWEATHER_API_KEY")
	githubUsername := "josephburgess"
	weatherLocation := "Buenosaires, AR"

	var mu sync.Mutex
	data := PageData{
		ProfileImage: "/static/images/profile.png",
		GithubURL:    "https://github.com/josephburgess",
		LinkedInURL:  "https://linkedin.com/in/josephburgessmba",
		Email:        "joe@joeburgess.dev",
		IsDarkMode:   false,
	}

	updateData(&data, githubUsername, weatherLocation, weatherApiKey, &mu)
	startBackgroundUpdater(&data, githubUsername, weatherLocation, weatherApiKey, &mu)

	lastModTime := time.Now()
	lastReloadTime := time.Now()

	go watchTemplate(tmplPath, &lastModTime, &tmpl)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/check-reload", func(w http.ResponseWriter, r *http.Request) {
		if lastModTime.After(lastReloadTime) {
			lastReloadTime = time.Now()
			w.Write([]byte("reload"))
		} else {
			w.Write([]byte("ok"))
		}
	})

	http.HandleFunc("/toggle-theme", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:     "theme",
			Value:    "light",
			Path:     "/",
			MaxAge:   86400 * 30,
			HttpOnly: true,
		}

		if data.IsDarkMode {
			cookie.Value = "light"
			data.IsDarkMode = false
		} else {
			cookie.Value = "dark"
			data.IsDarkMode = true
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/api/github-data", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		apiData := map[string]any{
			"repos":      data.GithubRepos,
			"activities": data.GitHubActivities,
			"updated":    data.LastUpdated,
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apiData)
	})

	http.HandleFunc("/update-data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		go updateData(&data, githubUsername, weatherLocation, weatherApiKey, &mu)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Data update triggered"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		cookie, err := r.Cookie("theme")
		if err == nil {
			data.IsDarkMode = (cookie.Value == "dark")
		}

		mu.Lock()
		localData := data
		mu.Unlock()

		err = tmpl.Execute(w, localData)
		if err != nil {
			log.Printf("Template execution error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func watchTemplate(tmplPath string, lastModTime *time.Time, tmplPtr **template.Template) {
	for {
		time.Sleep(1 * time.Second)

		info, err := os.Stat(tmplPath)
		if err != nil {
			continue
		}

		if info.ModTime().After(*lastModTime) {
			log.Println("Template changed, reloading...")

			newTmpl, err := template.ParseFiles(tmplPath)
			if err != nil {
				log.Printf("Error parsing template: %v", err)
				continue
			}

			*tmplPtr = newTmpl
			*lastModTime = info.ModTime()
			log.Println("Template reloaded successfully")
		}
	}
}
