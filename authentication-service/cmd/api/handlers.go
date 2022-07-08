package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Test(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Working properly from authentication"),
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) GetAuth(w http.ResponseWriter, r *http.Request) {

	user, err := app.Models.User.GetByEmail("admin@example.com")
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials "), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches("verysecret")
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials "), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// log.Fatal("from Authenticate in the auth service")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials "), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials "), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
