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
	"github.com/lardira/wicked-wit/internal/s3"
)

const (
	defaultEnvPath = "../.env"
)

func init() {
	err := godotenv.Load(defaultEnvPath)
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

	s3Config := &s3.S3Config{
		AccessKeyID:     os.Getenv("MINIO_SERVER_ACCESS_KEY"),
		SecretAccessKey: os.Getenv("MINIO_SERVER_SECRET_KEY"),
		Url:             os.Getenv("MINIO_URL"),
		Bucket:          os.Getenv("MINIO_BUCKET_NAME"),
	}

	err = s3.Init(s3Config)
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to connect to s3: %v\n", err))
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/games", handler.GameRouter())
	r.Mount("/cards", handler.CardRouter())
	r.Mount("/users", handler.UserRouter())

	log.Println("the server is running")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
