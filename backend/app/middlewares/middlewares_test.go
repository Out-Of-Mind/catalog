package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGreeter(t *testing.T) {
	w := httptest.NewRecorder()
	r := mux.NewRouter()
	Greeter("Hello").AddRoute(r)
	r.ServeHTTP(w, httptest.NewRequest("GET", "/greet/Hodor", nil))

	if w.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
	if w.Body.String() != "Hello Hodor!" {
		t.Error("Did not get expected greeting, got", w.Body.String())
	}
}