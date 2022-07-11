package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string `json:"action"`
	// Mail   MailPayload `json:"mail,omitempty"`
	Auth AuthPayload `json:"auth,omitempty"`
	// Log    LogPayload  `json:"log,omitempty"`
}
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Test(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Working properly from broker"),
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit the broker!",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Fatal("error in the readjson payload function")
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "getauth":
		app.GetAuth(w)
	default:
		app.errorJSON(w, errors.New("unknown action"))
		log.Fatal("wrong auth or other actions sent in the switch")
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	// create json and send it to the auth micro service
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		log.Fatal("error in MarshalIndent function")

		return
	}

	authServiceURL := fmt.Sprintf("http://%s/authenticate", "authentication-service:8083")

	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		log.Fatal("error in NewRequest function")
		return
	}
	request.Header.Set("Content-Type", "application/json")

	// we make sure we get back the cofrrect status code
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {

		app.errorJSON(w, err)
		log.Fatal("error in client Do function")
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetAuth(w http.ResponseWriter) {

	response, err := http.Get("http://authentication-service:8083/getauth")
	if err != nil {
		log.Fatal("error in posting function", err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	_ = app.writeJSON(w, http.StatusAccepted, payload)

}
