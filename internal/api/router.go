package api

import (
	"net/http"
	"time"

	"github.com/josephburgess/glogger"
	"github.com/josephburgess/joeburgess.dev/internal/api/handlers"
	"github.com/josephburgess/joeburgess.dev/internal/logging"
	"github.com/josephburgess/joeburgess.dev/internal/templates"
)

func Setup(tmplRenderer *templates.Renderer, dataUpdater *templates.DataUpdater) *http.Server {
	mux := http.NewServeMux()

	homeHandler := handlers.NewHomeHandler(tmplRenderer, dataUpdater)
	githubHandler := handlers.NewGithubHandler(dataUpdater)

	mux.HandleFunc("GET /{$}", homeHandler.HandleHome)
	mux.HandleFunc("POST /update-data", homeHandler.HandleUpdateData)
	mux.HandleFunc("GET /api/github-data", githubHandler.HandleGithubData)

	blog, err := glogger.New(glogger.Config{
		ContentDir: "content/posts",
		URLPrefix:  "/blog",
		Theme:      glogger.ThemeRosePine,
	})
	if err != nil {
		logging.Error("Failed to create blog", err)
	} else {
		prefix := blog.URLPrefix()
		mux.HandleFunc("GET "+prefix, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, prefix+"/", http.StatusMovedPermanently)
		})
		mux.Handle(prefix+"/", http.StripPrefix(prefix, blog.Handler()))
	}

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	handler := logging.Middleware(mux)

	return &http.Server{
		Addr:         ":8081",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
