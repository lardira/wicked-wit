package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/lardira/wicked-wit/handler"
	"github.com/lardira/wicked-wit/internal/db"
)

const (
	envPath = "../.env"
)

func init() {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	err := db.Init(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to connect to database: %v\n", err))
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/games", handler.GameRouter())

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
