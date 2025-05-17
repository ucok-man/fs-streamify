package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(app.withRecover)
	r.Use(app.withRecover)

	apiv1 := chi.NewRouter()
	apiv1.Route("/v1", func(r chi.Router) {
		/* -------------------------- Auth route -------------------------- */
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", app.signup)
			r.Post("/signin", app.signin)
			r.Post("/signout", app.signout)
		})
	})

	r.Mount("/api/v1", apiv1)
	return r
}
