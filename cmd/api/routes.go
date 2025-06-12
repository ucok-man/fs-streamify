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
	// r.Use(app.withCORS)

	apiv1 := chi.NewRouter()
	apiv1.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", app.signup)
			r.Post("/signin", app.signin)
			r.Post("/signout", app.signout)
			r.With(app.withAuthentication).Post("/onboarding", app.onboarding)
			r.With(app.withAuthentication).Get("/me", app.whoami)
		})
		r.Route("/users", func(r chi.Router) {
			r.Use(app.withAuthentication)

			r.Get("/{userId}", app.getUserById)

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

	if app.config.Env == "production" {
		staticDir := http.Dir("./build/ui")
		fs := http.FileServer(staticDir)

		// Serve static files and fallback to index.html for SPA routes
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			if _, err := staticDir.Open(r.URL.Path); err == nil {
				fs.ServeHTTP(w, r)
			} else {
				// Serve index.html for non-existent paths (e.g. /dashboard, /settings)
				http.ServeFile(w, r, "./build/ui/index.html")
			}
		})
	}

	r.Mount("/api/v1", apiv1)
	return r
}
