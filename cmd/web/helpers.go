package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	sdk "github.com/DeviaVir/servente-sdk"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) // don't report from this function but the originator

	if app.debug {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	response := &sdk.JSONResponse{
		Error:     true,
		Timestamp: time.Now(),
		Data: &sdk.JSONData{
			Message: http.StatusText(http.StatusInternalServerError),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	response := &sdk.JSONResponse{
		Error:     true,
		Timestamp: time.Now(),
		Data: &sdk.JSONData{
			Message: http.StatusText(status),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) renderJSON(w http.ResponseWriter, r *http.Request, d *sdk.JSONData) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := &sdk.JSONResponse{
		Error:     false,
		Timestamp: time.Now(),
		Data:      d,
	}

	json.NewEncoder(w).Encode(response)
}

func (app *application) parseDefaultRequest(r *http.Request) (string, string, error) {
	// do we have a pagination token?
	var nextPageToken string
	if nextPageToken == "" {
		nextPageToken = r.URL.Query().Get("page-token")
	}
	// try to find it in POST body
	if nextPageToken == "" {
		if err := r.ParseForm(); err != nil {
			return "", "", fmt.Errorf("could not parse form %v", err)
		}
		nextPageToken = r.PostForm.Get("page-token")
	}

	// do we have a MyCustomer ID?
	var myCustomer string
	if myCustomer == "" {
		myCustomer = r.URL.Query().Get("my-customer")
	}
	// try to find it in POST body
	if myCustomer == "" {
		if err := r.ParseForm(); err != nil {
			return "", "", fmt.Errorf("could not parse form %v", err)
		}
		myCustomer = r.PostForm.Get("my-customer")
	}
	if myCustomer == "" {
		myCustomer = "my_customer"
	}

	return nextPageToken, myCustomer, nil
}
