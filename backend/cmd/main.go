package main

import (
	"fmt"
	"net/http"

	"github.com/Rha02/carpool_app/handlers"
)

const port = ":8080"

func main() {
	fmt.Println("Starting App...")

	// Initialize the repository for handlers
	repo := handlers.NewRepo(0)
	handlers.NewHandlers(repo)

	router := routes()

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	server.ListenAndServe()
}
