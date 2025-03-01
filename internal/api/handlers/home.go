package handlers

import (
	"net/http"

	"github.com/josephburgess/joeburgess-dev/internal/templates"
)

type HomeHandler struct {
	renderer    *templates.Renderer
	dataUpdater *templates.DataUpdater
}

func NewHomeHandler(renderer *templates.Renderer, dataUpdater *templates.DataUpdater) *HomeHandler {
	return &HomeHandler{
		renderer:    renderer,
		dataUpdater: dataUpdater,
	}
}

func (h *HomeHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.HandleNotFound(w, r)
		return
	}

	data := h.dataUpdater.GetData()
	cookie, err := r.Cookie("theme")
	if err == nil {
		data.IsDarkMode = (cookie.Value == "dark")
	}

	html, err := h.renderer.RenderTemplate(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (h *HomeHandler) HandleUpdateData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	go h.dataUpdater.Update()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Data update triggered"))
}

func (h *HomeHandler) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Page not found"))
}
