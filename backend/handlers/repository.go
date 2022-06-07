package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Repo is a repository for handlers
var Repo *Repository

// Repository will contain any variables globally used by handlers
type Repository struct {
	dummy int
}

// NewRepo creates and returns a pointer to a new repository
func NewRepo(n int) *Repository {
	return &Repository{n}
}

// NewHandlers will set the global repo variable
func NewHandlers(r *Repository) {
	Repo = r
}

func respondJSON(rw http.ResponseWriter, msg interface{}, code int) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("X-Content-Type-Options", "nosniff")

	rw.WriteHeader(code)

	json.NewEncoder(rw).Encode(msg)
}

func (repo *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("This is Home")
}
