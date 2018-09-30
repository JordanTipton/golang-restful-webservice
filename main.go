package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jordantipton/golang-restful-webservice/apis"
	"github.com/jordantipton/golang-restful-webservice/app"
	"github.com/jordantipton/golang-restful-webservice/repositories"
	"github.com/jordantipton/golang-restful-webservice/services"
)

func main() {
	// Load configuration
	if err := app.LoadConfig("./config/config.json"); err != nil {
		panic(fmt.Errorf("Failed to load configuration: %s", err))
	}

	// Connect to database
	db, err := sql.Open("mysql", app.Config.DSN)
	if err != nil {
		panic(err)
	}
	r := buildRouter(db)
	http.ListenAndServe(":"+strconv.Itoa(app.Config.ServerPort), r)
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
