package handlers

import (
	"net/http"

	"github.com/Rha02/carpool_app/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *Repository) GetAllThreads(rw http.ResponseWriter, r *http.Request) {
	threads, err := m.DB.GetAllThreads()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(rw, threads, http.StatusOK)
}

func (m *Repository) GetUserThreads(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	threads, err := m.DB.GetUserThreads(vars["u_id"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(rw, threads, http.StatusOK)
}

func (m *Repository) GetThread(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	thread, err := m.DB.GetThreadByID(vars["id"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
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

	userID, err := primitive.ObjectIDFromHex(r.Form.Get("user_id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	title := r.Form.Get("title")
	body := r.Form.Get("body")

	t := models.Thread{
		UserID: userID,
		Title:  title,
		Body:   body,
	}

	err = m.DB.CreateThread(t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(rw, "", http.StatusCreated)
}

func (m *Repository) UpdateThread(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := r.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	id := vars["id"]
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	title := r.Form.Get("title")
	body := r.Form.Get("body")

	t := models.Thread{
		Title: title,
		Body:  body,
	}

	err = m.DB.UpdateThreadByID(id, t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(rw, "", http.StatusOK)
}

func (m *Repository) DeleteThread(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := m.DB.DeleteThreadByID(vars["id"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(rw, "", http.StatusOK)
}
