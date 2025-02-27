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

func fileWatcher(tmplPath string, lastModTime *time.Time, tmplPtr **template.Template) {
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

func main() {
	tmplPath := filepath.Join("templates", "index.html")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	info, err := os.Stat(tmplPath)
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}
	lastModTime := info.ModTime()

	go fileWatcher(tmplPath, &lastModTime, &tmpl)

	data := PageData{
		ProfileImage: "/static/images/profile.png",
		GithubURL:    "https://github.com/josephburgess",
		LinkedInURL:  "https://linkedin.com/in/josephburgessmba",
		Email:        "joe@joeburgess.dev",
		IsDarkMode:   false,
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/livereload.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(`
			(function() {
				const checkForReload = () => {
					fetch('/check-reload')
						.then(response => response.text())
						.then(data => {
							if (data === 'reload') {
								console.log('Reloading page...');
								window.location.reload();
							}
						})
						.catch(err => console.error('Error checking for reload:', err));
				};

				setInterval(checkForReload, 1000);
				console.log('Live reload enabled');
			})();
		`))
	})

	lastReloadTime := time.Now()

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
			MaxAge:   86400 * 30, // 30 days
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
	log.Println("Live reload enabled - changes to templates will auto-refresh the browser")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
