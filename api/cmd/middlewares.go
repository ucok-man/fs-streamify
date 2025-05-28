package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ucok-man/streamify-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func (app *application) withAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt-auth-token.streamify")
		if err != nil {
			app.errInvalidAuthenticationToken(w, r)
			return
		}
		var claim JWTClaim
		err = app.DecodeJwtToken(cookie.Value, &claim, app.config.JWT.AuthSecret)
		if err != nil {
			app.errInvalidAuthenticationToken(w, r)
			return
		}

		uid, err := bson.ObjectIDFromHex(claim.UserID)
		if err != nil {
			app.errInvalidAuthenticationToken(w, r)
			return
		}

		user, err := app.models.User.GetById(uid)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				app.errInvalidAuthenticationToken(w, r)
				return
			default:
				app.errInternalServer(w, r, err)
				return
			}
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
