package main

import (
	"net/http"

	"github.com/Rha02/carpool_app/handlers"
	"github.com/gorilla/mux"
)

func routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.Repo.Home).Methods("GET")

	router.HandleFunc("/users", handlers.Repo.UsersGetAll).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.Repo.UsersGet).Methods("GET")
	router.HandleFunc("/users", handlers.Repo.UsersPost).Methods("POST")
	router.HandleFunc("/users/{id}", handlers.Repo.UsersDelete).Methods("DELETE")

	return router
}
