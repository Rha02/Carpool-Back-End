package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rha02/Carpool-Back-End/config"
	"github.com/Rha02/Carpool-Back-End/dbrepo"
	"github.com/Rha02/Carpool-Back-End/driver"
	"github.com/Rha02/Carpool-Back-End/models"
	"github.com/gorilla/sessions"
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

func getSessionUser(s *sessions.Session) (*models.User, error) {
	var user models.User

	user, ok := s.Values["user"].(models.User)
	if !ok {
		return nil, fmt.Errorf("error: could not get user from session")
	}

	return &user, nil
}
