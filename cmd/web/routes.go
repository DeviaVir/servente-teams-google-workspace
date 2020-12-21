package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.authenticate)

	mux := pat.New()

	// WARNING: do not change the below routes, servente expects these to be available for
	// all providers.
	mux.Get("/api", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.apiList))
	mux.Get("/api/v1/teams/list", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.teamsList))
	mux.Get("/api/v1/teams/membership", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.teamsMembership))
	mux.Get("/api/v1/teams/membership/:member", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.teamsMembership))
	// WARNING: end

	mux.Get("/ping", http.HandlerFunc(ping))

	return standardMiddleware.Then(mux)
}
