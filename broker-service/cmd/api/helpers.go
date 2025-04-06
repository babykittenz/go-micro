package main

import (
	toolkit "github.com/babykittenz/api-micro-util"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	tools := toolkit.Tools{}
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
		Success: true,
	}

	_ = tools.WriteJSON(w, http.StatusOK, payload)
}
