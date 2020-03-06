package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandleThings(t *testing.T) {
	req, err := http.NewRequest("GET", "/things", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	s := newServer()
	handler := http.HandlerFunc(s.handleThings())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error("Failed")
	}
}

func TestPostNotAllowedHandleThings(t *testing.T) {
	req, err := http.NewRequest("POST", "/things", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	s := newServer()
	handler := http.HandlerFunc(s.handleThings())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Error("Failed")
	}
}

func TestRecoverMiddlewareHandlesPanic(t *testing.T) {
	req, err := http.NewRequest("GET", "/things", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(panicHandler())
	wrapped := recoverMiddleware(handler)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Error("should have got a 500")
	}
}

func panicHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("I blew up")
	}
}
