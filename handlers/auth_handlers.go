package handlers

import (
	"net/http"

	"github.com/Rha02/Carpool-Back-End/models"
	"github.com/Rha02/Carpool-Back-End/utils"
	"github.com/gorilla/csrf"
)

// GetToken responds with a CSRF token
func (m *Repository) GetToken(rw http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	rw.Header().Set("X-CSRF-Token", csrfToken)

	respondJSON(rw, "", http.StatusOK)
}

func (m *Repository) Login(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	phone := r.Form.Get("phone")
	password := r.Form.Get("password")

	u, err := m.DB.Authenticate(phone, password)
	if err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	u.Password = ""

	session.Values["user"] = u
	session.Save(r, rw)

	respondJSON(rw, "", http.StatusOK)
}

func (m *Repository) Logout(rw http.ResponseWriter, r *http.Request) {
	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, "user")
	session.Save(r, rw)

	respondJSON(rw, "", http.StatusOK)
}

// Register registers a new user and authenticates them
func (m *Repository) Register(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	phone := r.Form.Get("phone")
	name := r.Form.Get("name")
	password := r.Form.Get("password")

	u := models.User{
		Phone:    phone,
		Name:     name,
		Password: password,
	}

	user, err := m.DB.RegisterUser(u)
	if err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Do not store hashed password in session
	user.Password = ""

	session.Values["user"] = *user
	session.Save(r, rw)

	respondJSON(rw, "", http.StatusOK)
}

// CheckAuth
func (m *Repository) CheckAuth(rw http.ResponseWriter, r *http.Request) {
	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	user, ok := session.Values["user"]
	if ok {
		respondJSON(rw, user, http.StatusOK)
		return
	}

	http.Error(rw, "Forbidden", http.StatusForbidden)
}
