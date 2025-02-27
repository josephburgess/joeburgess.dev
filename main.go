package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type PageData struct {
	ProfileImage string
	GithubURL    string
	LinkedInURL  string
	Email        string
	IsDarkMode   bool
}

func main() {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	data := PageData{
		ProfileImage: "/static/images/profile.png",
		GithubURL:    "https://github.com/josephburgess",
		LinkedInURL:  "https://linkedin.com/in/josephburgessmba",
		Email:        "joe@joeburgess.dev",
		IsDarkMode:   false,
	}

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("theme")
		if err == nil {
			data.IsDarkMode = (cookie.Value == "dark")
		}

		err = tmpl.Execute(w, data)
		if err != nil {
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
