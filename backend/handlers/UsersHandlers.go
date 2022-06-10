package handlers

import (
	"fmt"
	"net/http"

	"github.com/Rha02/carpool_app/models"
	"github.com/gorilla/mux"
)

// UsersGetAll returns all users
func (m *Repository) GetAllUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	users := m.DB.GetAllUsers()

	respondJSON(rw, users, http.StatusOK)
}

// GetUsers returns a user by specified id
func (m *Repository) GetUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := m.DB.GetUserByID(vars["id"])
	if err != nil {
		http.Error(rw, fmt.Sprintf("error: could not find user by id %s", vars["id"]), http.StatusNotFound)
		return
	}

	// TODO: Return helpful json error response when user is not found

	respondJSON(rw, *user, http.StatusOK)
}

func (m *Repository) PostUser(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(rw, "error: could not parse request form", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")

	// TODO: Add email validation

	name := r.Form.Get("name")

	user := models.User{
		Email: email,
		Name:  name,
	}

	err = m.DB.CreateUser(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusSeeOther)
	}

	respondJSON(rw, "", http.StatusOK)
}

func (m *Repository) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := m.DB.DeleteUserByID(vars["id"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	respondJSON(rw, "", http.StatusOK)
}
