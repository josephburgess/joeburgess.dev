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

	if rr.Code != http.StatusNotFound {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusNotFound)
	}
}
