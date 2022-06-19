package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rha02/carpool_app/config"
	"github.com/Rha02/carpool_app/dbrepo"
	"github.com/Rha02/carpool_app/driver"
)

// Repo is a repository for handlers
var Repo *Repository

// Repository will contain any variables globally used by handlers
type Repository struct {
	App *config.AppConfig
	DB  dbrepo.DatabaseRepository
}

// NewRepo creates and returns a pointer to a new repository
func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewDatabaseRepo(db),
	}
}

func NewTestRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewTestingRepo(),
	}
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
