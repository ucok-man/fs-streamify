package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.NotFound(app.errNotFound)
	r.MethodNotAllowed(app.errMethodNotAllowed)

	r.Use(app.withRecover)
	r.Use(app.withCORS)

	apiv1 := chi.NewRouter()
	apiv1.Group(func(r chi.Router) {
		/* -------------------------- Auth route -------------------------- */
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", app.signup)
			r.Post("/signin", app.signin)
			r.Post("/signout", app.signout)
			r.With(app.withAuthentication).Post("/onboarding", app.onboarding)
			r.With(app.withAuthentication).Get("/me", app.whoami)
		})
		r.Route("/users", func(r chi.Router) {
			r.Use(app.withAuthentication)

			r.Get("/recommended", app.recommended)
			r.Get("/friends-with-me", app.myfriend)

			r.Route("/friends-request", func(r chi.Router) {
				r.Post("/create/{recipientId}", app.requestFriend)
				r.Post("/accept/{friendRequestId}", app.acceptFriend)
				r.Get("/from", app.getAllFromFriendRequest)
				r.Get("/send", app.getAllSendFriendRequest)
			})
		})

		r.Route("/chat", func(r chi.Router) {
			r.Use(app.withAuthentication)
			r.Get("/token", app.getStreamToken)
		})
	})

	r.Mount("/api/v1", apiv1)
	return r
}
