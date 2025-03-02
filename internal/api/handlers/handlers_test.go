package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleNotFound(t *testing.T) {
	handler := &HomeHandler{}

	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.HandleNotFound(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	expected := "404 - Page not found"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandleUpdateDataMethodNotAllowed(t *testing.T) {
	handler := &HomeHandler{}

	req, err := http.NewRequest("GET", "/update", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.HandleUpdateData(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestWeatherHandlerMissingLocation(t *testing.T) {
	handler := &WeatherHandler{}

	req, err := http.NewRequest("GET", "/api/weather", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.HandleWeatherData(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "Location parameter is required"
	if rr.Body.String() != expected+"\n" {
		t.Errorf("handler returned unexpected body: got %q want %q",
			rr.Body.String(), expected+"\n")
	}
}
