package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/vishal-rfx/snippetbox/ui"
)

// routes method returns a servemux containing application routes.
func (app *application) routes() http.Handler {

	// Use the http.NewServeMux() function to initialize a new servermux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Use the http.FileServerFS() function to create a HTTP handler which
	// serves the embedded files in ui.Files. It's important to note that our static files are 
	// contained in the static folder of the ui.Files embedded filesystem. So, for example, our CSS stylesheet is located at 
	// "static/css/main.css".
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	
	// Add a new GET /ping route.
	mux.HandleFunc("GET /ping", ping)

	// Create a new middleware chain containing the middleware specific to our dynamic application routes.
	// For now, this chain will only contain the LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

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