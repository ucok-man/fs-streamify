package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.
		Err(err).
		CallerSkipFrame(2).
		Str("request_method", r.Method).
		Str("request_url", r.URL.String()).Send()
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) errInternalServer(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, errors.WithStack(err))

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, 500, message)
}

func (app *application) errNotFound(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) errMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) errBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) errFailedValidation(w http.ResponseWriter, r *http.Request, errors map[string][]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) errEditConflict(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) errRateLimitExceeded(w http.ResponseWriter, r *http.Request) {
	message := "rate limited exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (app *application) errInvalidCredentials(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) errInvalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) errAuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) errNotPermitted(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
