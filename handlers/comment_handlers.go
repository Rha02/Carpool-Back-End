package handlers

import (
	"net/http"

	"github.com/Rha02/Carpool-Back-End/models"
	"github.com/Rha02/Carpool-Back-End/utils"
	"github.com/gorilla/mux"
)

func (m *Repository) PostComment(rw http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	threadID := vars["t_id"]

	comment := models.Comment{
		UserID: u.ID,
		Body:   r.Form.Get("body"),
	}

	if err := m.DB.CreateComment(threadID, comment); err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	respondJSON(rw, "", http.StatusCreated)
}

func (m *Repository) UpdateComment(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := m.App.CookieStore.Get(r, "session_id")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := getSessionUser(session)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	comment, err := m.DB.GetComment(vars["t_id"], vars["c_id"])
	if err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	if comment.UserID != u.ID {
		http.Error(rw, "error: client is forbidden from accessing this resource", http.StatusForbidden)
		return
	}

	comment.Body = r.Form.Get("body")

	if err = m.DB.UpdateComment(vars["t_id"], vars["c_id"], *comment); err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}

	respondJSON(rw, "", http.StatusOK)
}

func (m *Repository) DeleteComment(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	if err := m.DB.DeleteComment(vars["t_id"], vars["c_id"]); err != nil {
		dberr := err.(*utils.DBError)
		http.Error(rw, dberr.Error(), dberr.StatusCode())
		return
	}
	respondJSON(rw, "", http.StatusOK)
}
