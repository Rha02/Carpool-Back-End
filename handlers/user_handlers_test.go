package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// Tests
func TestGetAllUsers(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.GetAllUsers)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Unexpected status code. expected %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestGetUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/id", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "123",
	})

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.GetUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Unexpected status code: expected %d, got %d", http.StatusOK, rr.Code)
	}

	req, _ = http.NewRequest("GET", "/users/id", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "error",
	})

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.GetUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Unexpected status code: expected %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/users/id", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "123",
	})

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.DeleteUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Unexpected status code. expected %d, got %d", http.StatusOK, rr.Code)
	}

	req, _ = http.NewRequest("DELETE", "/users/id", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "error",
	})

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.DeleteUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Unexpected status code. expected %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	params := strings.NewReader(url.Values{
		"email": {"heisenberg@example.loc"},
		"name":  {"Heisenberg"},
	}.Encode())

	req, _ := http.NewRequest("PUT", "/users/id", params)
	req = mux.SetURLVars(req, map[string]string{
		"id": "123",
	})

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.UpdateUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Unexpected status code. expected %d, got %d", http.StatusOK, rr.Code)
	}

	params = strings.NewReader(url.Values{
		"email": {"heisenberg@example.loc"},
		"name":  {"Heisenberg"},
	}.Encode())

	req, _ = http.NewRequest("PUT", "/users/id", params)
	req = mux.SetURLVars(req, map[string]string{
		"id": "error",
	})

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.UpdateUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotAcceptable {
		t.Errorf("Unexpected status code. expected %d, got %d", http.StatusNotAcceptable, rr.Code)
	}
}
