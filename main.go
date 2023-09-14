package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpcontroller "go-shortener/controllers/http"
	"go-shortener/repositories/postgres"
	"net/http"
)

// for now hard coded
const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "root"
	dbname   = "go_shortener"
)

func main() {
	repo, err := postgres.NewPostgresURLRepository(
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname),
	)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Put("/u", func(w http.ResponseWriter, r *http.Request) {
		httpcontroller.NewURLController("http://localhost:3000", w, r, repo).Put()
	})

	http.ListenAndServe(":3000", r)
}
