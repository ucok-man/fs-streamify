package main

import "net/http"

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func (app *application) signin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func (app *application) signout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
