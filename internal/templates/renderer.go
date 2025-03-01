package templates

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/josephburgess/joeburgess-dev/internal/models"
)

type PageData struct {
	ProfileImage     string
	GithubURL        string
	LinkedInURL      string
	BreezeURL        string
	Email            string
	IsDarkMode       bool
	GithubRepos      []models.Repository
	GitHubActivities []models.Activity
	LastUpdated      string
	Weather          *models.WeatherData
	mu               sync.RWMutex
}

type Renderer struct {
	tmpl        *template.Template
	lastModTime time.Time
	tmplPath    string
	mu          sync.RWMutex
}

func NewRenderer() *Renderer {
	tmplPath := filepath.Join("templates", "index.html")

	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"formatDate": formatDate,
		"timeSince":  timeSince,
		"toLower":    strings.ToLower,
	}).ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	renderer := &Renderer{
		tmpl:        tmpl,
		lastModTime: time.Now(),
		tmplPath:    tmplPath,
	}

	go renderer.watchTemplate()

	return renderer
}

func (r *Renderer) RenderTemplate(data *PageData) (template.HTML, error) {
	r.mu.RLock()
	tmpl := r.tmpl
	r.mu.RUnlock()

	return executeToHTML(tmpl, data)
}

func executeToHTML(tmpl *template.Template, data any) (template.HTML, error) {
	var sb strings.Builder
	if err := tmpl.Execute(&sb, data); err != nil {
		return "", err
	}
	return template.HTML(sb.String()), nil
}

func (r *Renderer) watchTemplate() {
	for {
		time.Sleep(1 * time.Second)

		info, err := os.Stat(r.tmplPath)
		if err != nil {
			continue
		}

		if info.ModTime().After(r.lastModTime) {
			log.Println("Template changed, reloading...")

			newTmpl, err := template.New("index.html").Funcs(template.FuncMap{
				"formatDate": formatDate,
				"timeSince":  timeSince,
				"toLower":    strings.ToLower,
			}).ParseFiles(r.tmplPath)
			if err != nil {
				log.Printf("Error parsing template: %v", err)
				continue
			}

			r.mu.Lock()
			r.tmpl = newTmpl
			r.lastModTime = info.ModTime()
			r.mu.Unlock()

			log.Println("Template reloaded successfully")
		}
	}
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
