package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lardira/wicked-wit/handler"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/games", handler.GameRouter())

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
