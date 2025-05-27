package main

import (
	"net/http"
	"time"
)

func (app *application) getStreamToken(w http.ResponseWriter, r *http.Request) {
	currentUser := app.contextGetUser(r)
	token, err := app.stream.CreateToken(currentUser.ID.Hex(), time.Time{}) // zero time expired
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"token": token}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}
