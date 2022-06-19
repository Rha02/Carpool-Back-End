package handlers

import (
	"net/http"

	"github.com/Rha02/carpool_app/models"
)

func (m *Repository) Login(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	err = m.DB.Authenticate(email, password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Save(r, rw)

	respondJSON(rw, "", http.StatusOK)
}

func (m *Repository) Logout(rw http.ResponseWriter, r *http.Request) {
	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = false
	session.Save(r, rw)

	respondJSON(rw, "", http.StatusOK)
}

func (m *Repository) Register(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	email := r.Form.Get("email")
	name := r.Form.Get("name")
	password := r.Form.Get("password")

	u := models.User{
		Email:    email,
		Name:     name,
		Password: password,
	}

	err = m.DB.RegisterUser(u)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Save(r, rw)

	respondJSON(rw, "", http.StatusOK)
}

func (m *Repository) CheckAuth(rw http.ResponseWriter, r *http.Request) {
	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	auth, ok := session.Values["authenticated"]
	if ok && auth != false {
		rw.Write([]byte("Authentication Successful!"))
		return
	}

	http.Error(rw, "Forbidden", http.StatusForbidden)
}
