package main

import (
	"net/http"
	"time"

	response "github.com/ucok-man/streamify-api/cmd/responses"
)

func (app *application) getStreamToken(w http.ResponseWriter, r *http.Request) {
	currentUser := app.contextGetUser(r)
	token, err := app.stream.CreateToken(currentUser.ID.Hex(), time.Time{}) // zero time expired
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}
	var payload response.TokenResponse
	payload.Value = token

	err = app.writeJSON(w, http.StatusCreated, envelope{"token": payload}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}
