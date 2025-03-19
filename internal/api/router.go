package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/josephburgess/glogger"
	"github.com/josephburgess/joeburgess.dev/internal/api/handlers"
	"github.com/josephburgess/joeburgess.dev/internal/logging"
	"github.com/josephburgess/joeburgess.dev/internal/templates"
)

func Setup(tmplRenderer *templates.Renderer, dataUpdater *templates.DataUpdater) *http.Server {
	router := mux.NewRouter()

	homeHandler := handlers.NewHomeHandler(tmplRenderer, dataUpdater)
	githubHandler := handlers.NewGithubHandler(dataUpdater)

	router.Use(logging.Middleware)

	router.HandleFunc("/", homeHandler.HandleHome).Methods("GET")
	router.HandleFunc("/update-data", homeHandler.HandleUpdateData).Methods("POST")
	router.HandleFunc("/api/github-data", githubHandler.HandleGithubData).Methods("GET")

	blog, err := glogger.New(glogger.Config{
		ContentDir: "content/posts",
		URLPrefix:  "/blog",
	})
	if err != nil {
		logging.Error("Failed to create blog", err)
	} else {
		blog.RegisterHandlers(router)
	}

	fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	router.NotFoundHandler = http.HandlerFunc(homeHandler.HandleNotFound)

	return &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
