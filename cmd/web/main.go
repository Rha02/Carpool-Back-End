package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/Rha02/Carpool-Back-End/config"
	"github.com/Rha02/Carpool-Back-End/driver"
	"github.com/Rha02/Carpool-Back-End/handlers"
	"github.com/Rha02/Carpool-Back-End/models"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

const port = ":8080"

var app config.AppConfig

func main() {
	fmt.Println("Starting App...")

	godotenv.Load("../../.env")

	uri := fmt.Sprintf("mongodb://%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db, err := driver.ConnectMongoDB(uri, os.Getenv("DB_NAME"))
	if err != nil {
		panic(err)
	}

	// Register User model for use by gorilla/sessions
	gob.Register(models.User{})

	store := sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))

	app.CookieStore = store

	// Initialize the repository for handlers
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	auth_key := os.Getenv("AUTH_KEY")
	router := routes(auth_key)

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Printf("App Running on Port %s\n", port)
	server.ListenAndServe()
}
