package main

import (
	"fmt"
	"logger-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Test(w http.ResponseWriter, r *http.Request) {
	// log.Fatal("working well from Test function in handlers")
	payload := JSONPayload{
		Name: "working",
		Data: fmt.Sprintf("Working properly from logger service"),
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// log.Fatal("working well from WriteLog function in handlers")

	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert the data
	// err := app.logEvent(requestPayload.Name, requestPayload.Data)
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	// create the response we'll send back as JSON

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}
	// write the response back as JSON
	_ = app.writeJSON(w, http.StatusAccepted, resp)
}
