package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
)

// serverError helper writes a log entry at Error level (including the request method and request URI as attributes),
// then sends a generic 500 Internal server error response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error){
	var (
		method = r.Method
		uri = r.RequestURI
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description to the user. We'll 
// use this later in the book to send responses like 400 "Bad Request" when there's a problem with the
// request that user sent.
func (app *application) clientError(w http.ResponseWriter,r *http.Request, status int){
	http.Error(w, http.StatusText(status), status)
} 


func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData){
	// Retrieve the appropriate template set from the cache based on the page name
	// like ('home.html.tmpl'). If no entry exists in the cache with the provided name, then create a new error and call
	// the serverError() helper method
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the http.ResponseWriter. If there's an error, call our
	// serverError() helper and then return
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}


	// Write out the provided HTTP status code (200 - OK, 400 - Bad Request, etc..)
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)

}

// newTemplateData returns a templateData struct initialized with the current year. Note that we're not 
func (app *application) newTemplateData(r *http.Request) templateData{
	return templateData{
		CurrentYear: time.Now().Year(),
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken: nosurf.Token(r),
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}