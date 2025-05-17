package main

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) withRecover(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

func (app *application) withCORS(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   app.config.Cors.Origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           60, // in seconds
	})(next)

}
