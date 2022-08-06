package handlers

import (
	"net/http"

	"github.com/Rha02/Carpool-Back-End/models"
	"github.com/Rha02/Carpool-Back-End/utils"
	"github.com/gorilla/mux"
)

func (m *Repository) GetAllThreads(rw http.ResponseWriter, r *http.Request) {
	threads, err := m.DB.GetAllThreads()
	if err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	respondJSON(rw, threads, http.StatusOK)
}

func (m *Repository) GetUserThreads(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	threads, err := m.DB.GetUserThreads(vars["u_id"])
	if err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	respondJSON(rw, threads, http.StatusOK)
}

func (m *Repository) GetThread(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	thread, err := m.DB.GetThreadByID(vars["id"])
	if err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	respondJSON(rw, thread, http.StatusOK)
}

func (m *Repository) PostThread(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := getSessionUser(session)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	title := r.Form.Get("title")
	body := r.Form.Get("body")

	t := models.Thread{
		UserID:   user.ID,
		Title:    title,
		Body:     body,
		Comments: []models.Comment{},
	}

	if err = m.DB.CreateThread(t); err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	respondJSON(rw, "", http.StatusCreated)
}

func (m *Repository) UpdateThread(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
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

	t, err := m.DB.GetThreadByID(vars["id"])
	if err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	if t.UserID != u.ID {
		http.Error(rw, "error: client has no access to this resource", http.StatusForbidden)
		return
	}

	t.Title = r.Form.Get("title")
	t.Body = r.Form.Get("body")

	if err = m.DB.UpdateThreadByID(vars["id"], *t); err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	respondJSON(rw, "", http.StatusOK)
}

func (m *Repository) DeleteThread(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

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

	t, err := m.DB.GetThreadByID(vars["id"])
	if err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	if u.ID != t.UserID {
		http.Error(rw, "error: client has no access to this resource", http.StatusForbidden)
		return
	}

	if err = m.DB.DeleteThreadByID(vars["id"]); err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	respondJSON(rw, "", http.StatusOK)
}
