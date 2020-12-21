package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	sdk "github.com/DeviaVir/servente-sdk"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
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

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, statefulStorage string, errorLog *log.Logger) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := filepath.Join(statefulStorage, "servente-teams-google-workspace-token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config, errorLog)
		saveToken(tokFile, tok, errorLog)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config, errorLog *log.Logger) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		errorLog.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		errorLog.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token, errorLog *log.Logger) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		errorLog.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
