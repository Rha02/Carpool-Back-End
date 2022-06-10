package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Rha02/carpool_app/driver"
	"github.com/Rha02/carpool_app/handlers"
	"github.com/joho/godotenv"
)

const port = ":8080"

func main() {
	fmt.Println("Starting App...")

	godotenv.Load("../.env")

	uri := fmt.Sprintf("mongodb://%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db, err := driver.ConnectMongoDB(uri, os.Getenv("DB_NAME"))
	if err != nil {
		panic(err)
	}

	// Initialize the repository for handlers
	repo := handlers.NewRepo(db)
	handlers.NewHandlers(repo)

	router := routes()

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Printf("App Running on Port %s\n", os.Getenv("DB_PORT"))
	server.ListenAndServe()
}
