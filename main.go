package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jordantipton/golang-restful-webservice/apis"
	"github.com/jordantipton/golang-restful-webservice/repositories"
	"github.com/jordantipton/golang-restful-webservice/services"
)

func main() {
	// Connect to database
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}
	r := buildRouter(db)
	http.ListenAndServe(":8080", r)
}

func buildRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	// Middleware stack
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Register Controllers
	usersRepository := &repositories.UsersRepository{DB: db}
	usersService := &services.UsersService{UsersPersister: usersRepository}
	apis.RegisterUsersResource(r, usersService)
	return r
}
