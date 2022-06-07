package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Rha02/carpool_app/models"
	"github.com/gorilla/mux"
)

// UsersGetAll returns all users
func (repo *Repository) UsersGetAll(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	users := models.GetExampleUsers()

	respondJSON(rw, users, http.StatusOK)
}

// UsersGet returns a user by specified id
func (repo *Repository) UsersGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "error: could not parse user id", http.StatusBadRequest)
		return
	}

	user := models.GetUserByID(id)

	// TODO: Return helpful json error response when user is not found
	if user == nil {
		http.Error(rw, fmt.Sprintf("error: user by id %d is not found", id), http.StatusNotFound)
		return
	}

	respondJSON(rw, *user, http.StatusOK)
}

func (repo *Repository) UsersPost(rw http.ResponseWriter, r *http.Request) {
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

	models.InsertUser(user)

	respondJSON(rw, "", http.StatusOK)
}

func (repo *Repository) UsersDelete(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "error: could not parse user id", http.StatusBadRequest)
		return
	}

	err = models.DeleteUser(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	respondJSON(rw, "", http.StatusOK)
}
