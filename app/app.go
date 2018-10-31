package app

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql" // Register mysql driver

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jordantipton/golang-restful-webservice/apis"
	"github.com/jordantipton/golang-restful-webservice/repositories"
	"github.com/jordantipton/golang-restful-webservice/services"
)

// App struct
type App struct {
	Router *chi.Mux
}

// Initialize app and construct router
func (a *App) Initialize(dsn string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	a.Router = buildRouter(db)
}

// Run app
func (a *App) Run(addr string) {
	http.ListenAndServe(addr, a.Router)
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
