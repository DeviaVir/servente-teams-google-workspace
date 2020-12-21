package main

import (
	"net/http"

	sdk "github.com/DeviaVir/servente-sdk"
)

// ping: healthcheck endpoint
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) apiList(w http.ResponseWriter, r *http.Request) {
	supportedAPIs := []string{"v1"}

	var apiVersions []sdk.JSONAPI
	for _, api := range supportedAPIs {
		apiVersions = append(apiVersions, sdk.JSONAPI{
			ID: api,
		})
	}

	response := &sdk.JSONData{
		APIVersions: apiVersions,
	}

	app.renderJSON(w, r, response)
}

// teamsList retrieve a list of teams (groups) and return in a JSON
func (app *application) teamsList(w http.ResponseWriter, r *http.Request) {
	nextPageToken, myCustomer, err := app.parseDefaultRequest(r)
	if err != nil {
		app.errorLog.Printf("%s", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	groupsResponse, err := app.googleClient.Groups.List().Customer(myCustomer).OrderBy("email").PageToken(nextPageToken).Do()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var teams []sdk.JSONTeam
	for _, team := range groupsResponse.Groups {
		tS := sdk.JSONTeam{
			Name:         team.Name,
			ID:           team.Id,
			Email:        team.Email,
			MembersCount: team.DirectMembersCount,
		}
		teams = append(teams, tS)
	}

	response := &sdk.JSONData{
		Teams:         teams,
		NextPageToken: groupsResponse.NextPageToken,
	}

	app.renderJSON(w, r, response)
}

// teamsMembership returns a list of teams (groups) an email is part of, return
// in a JSON
func (app *application) teamsMembership(w http.ResponseWriter, r *http.Request) {
	var member string
	if member == "" {
		member = r.URL.Query().Get(":member")
	}
	if member == "" {
		member = r.URL.Query().Get("member")
	}
	// try to find it in POST body
	if member == "" {
		if err := r.ParseForm(); err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		member = r.PostForm.Get("member")
	}

	if member == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	nextPageToken, _, err := app.parseDefaultRequest(r)
	if err != nil {
		app.errorLog.Printf("%s", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	groupsResponse, err := app.googleClient.Groups.List().UserKey(member).OrderBy("email").PageToken(nextPageToken).Do()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var teams []sdk.JSONTeam
	for _, team := range groupsResponse.Groups {
		tS := sdk.JSONTeam{
			Name:         team.Name,
			ID:           team.Id,
			Email:        team.Email,
			MembersCount: team.DirectMembersCount,
		}
		teams = append(teams, tS)
	}

	response := &sdk.JSONData{
		Teams:         teams,
		NextPageToken: groupsResponse.NextPageToken,
	}

	app.renderJSON(w, r, response)
}
