package handlers

import (
	"net/http"

	"github.com/Rha02/carpool_app/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UsersGetAll returns all users
func (m *Repository) GetAllUsers(rw http.ResponseWriter, r *http.Request) {
	users, err := m.DB.GetAllUsers()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(rw, users, http.StatusOK)
}

// GetUsers returns a user by specified id
func (m *Repository) GetUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := m.DB.GetUserByID(vars["id"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Return helpful json error response when user is not found

	respondJSON(rw, *user, http.StatusOK)
}

func (m *Repository) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	// TODO: Figure out how to handle user deletion.
	respondJSON(rw, "", http.StatusNotImplemented)
}

func (m *Repository) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusSeeOther)
		return
	}

	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := getSessionUser(session)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if u.ID != objectID {
		http.Error(rw, "error: user is forbidden from accessing this resource", http.StatusForbidden)
		return
	}

	email := r.Form.Get("email")
	name := r.Form.Get("name")

	updatedUser := models.User{
		ID:    objectID,
		Email: email,
		Name:  name,
	}

	err = m.DB.UpdateUserByID(id, updatedUser)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotAcceptable)
		return
	}

	// Update user value stored in session
	session.Values["user"] = updatedUser
	session.Save(r, rw)

	respondJSON(rw, "", http.StatusOK)
}
