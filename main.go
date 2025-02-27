package main

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	ProfileImage string
	GithubURL    string
	LinkedInURL  string
	Email        string
	IsDarkMode   bool
}

func main() {
	tmpl, err := template.ParseFiles("templates/index.html")
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

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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
