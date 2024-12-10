package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes method returns a servemux containing application routes.
func (app *application) routes() http.Handler {

	// Use the http.NewServeMux() function to initialize a new servermux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() to register the file server as the handler for all URL paths that start with 
	// "/static/" for matching paths, we strip the "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware specific to our dynamic application routes.
	// For now, this chain will only contain the LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	// Update these routes to use the new dynamic middleware chain followed by the
	// appropriate handler function. Because alice ThenFunc() method returns a http.Handler (rather than a http.HandlerFunc)
	// we also need to switch to registering the route using the mux.Handle() method.

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create/{$}", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create/{$}",protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))


	// Create a middleware chain containing our 'standard' middleware which will be used for every request our 
	// application receives
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)

}