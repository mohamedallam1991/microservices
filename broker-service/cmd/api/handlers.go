package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit the broker!",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
	// out, _ := json.MarshalIndent(payload, "", "\t")
	// w.Header().Set("content-Type", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out)
}
