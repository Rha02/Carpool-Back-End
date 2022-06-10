package main

import (
	"net/http"

	"github.com/Rha02/carpool_app/handlers"
	"github.com/gorilla/mux"
)

func routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.Repo.Home).Methods("GET")

	router.HandleFunc("/users", handlers.Repo.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.Repo.GetUser).Methods("GET")
	router.HandleFunc("/users", handlers.Repo.PostUser).Methods("POST")
	router.HandleFunc("/users/{id}", handlers.Repo.DeleteUser).Methods("DELETE")

	return router
}
