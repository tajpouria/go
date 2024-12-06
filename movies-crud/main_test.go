package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetMovies() {
	movies = []Movie{}
}

func TestMoviesIndex(t *testing.T) {
	resetMovies()

	req, err := http.NewRequest("GET", "/movies", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	moviesIndex(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expected := "[]\n" // Empty movies list
	if w.Body.String() != expected {
		t.Errorf("Expected body to be %s, got %s", expected, w.Body.String())
	}
}

func TestMoviesCreate(t *testing.T) {
	resetMovies()

	movie := Movie{Title: "The Godfather", Year: 1972}

	movieJson, err := json.Marshal(movie)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/movies/", bytes.NewBuffer(movieJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	moviesChange(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	expected := `{"id":1,"title":"The Godfather","year":1972}` + "\n"
	if w.Body.String() != expected {
		t.Errorf("Expected body to be %s, got %s", expected, w.Body.String())
	}
}

func TestMoviesCreateBadRequest(t *testing.T) {
	resetMovies()

	req, err := http.NewRequest("POST", "/movies/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	moviesChange(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMoviesCreateInvalidJson(t *testing.T) {
	resetMovies()

	req, err := http.NewRequest("POST", "/movies/", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	moviesChange(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMoviesDelete(t *testing.T) {
	resetMovies()

	movie := Movie{Title: "The Godfather", Year: 1972}
	movieJson, err := json.Marshal(movie)
	if err != nil {
		t.Fatal(err)
	}

	reqCreate, err := http.NewRequest("POST", "/movies/", bytes.NewBuffer(movieJson))
	if err != nil {
		t.Fatal(err)
	}
	reqCreate.Header.Set("Content-Type", "application/json")

	wCreate := httptest.NewRecorder()
	moviesChange(wCreate, reqCreate)

	if wCreate.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, wCreate.Code)
	}

	reqDelete, err := http.NewRequest("DELETE", "/movies/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	wDelete := httptest.NewRecorder()
	moviesChange(wDelete, reqDelete)

	if wDelete.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, wDelete.Code)
	}

	reqGet, err := http.NewRequest("GET", "/movies", nil)
	if err != nil {
		t.Fatal(err)
	}

	wGet := httptest.NewRecorder()
	moviesIndex(wGet, reqGet)

	expected := "[]\n"
	if wGet.Body.String() != expected {
		t.Errorf("Expected body to be %s, got %s", expected, wGet.Body.String())
	}
}
