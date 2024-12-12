package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/vishal-rfx/snippetbox/internal/models"
	"github.com/vishal-rfx/snippetbox/ui"
)

// humanDate function which returns a nicely formatter string representation of a time.Time object.
func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is essentially 
// a string-keyed map which acts as a lookup between the names of our custom template functions and the 
// functions themselves

var functions = template.FuncMap{
	"humanDate": humanDate,
}


// Create a template cache
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// use the fs.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl"
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Extract the filename from the full filepath and assign it to the name variable
		name := filepath.Base(page)
		patterns := []string {
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}


		// The template.FuncMap must be registered with the template set before we call the ParseFiles() method. This
		// means we have to use the template.New() to create an empty template set, use the Funcs method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		
		cache[name] = ts
	}
	return cache, nil
}




// Define a templateData type to act as the holding structure for any dynamic data that we want
// to pass to our HTML templates
// At the moment it only contains one field, but we'll add more to it as the build progresses

type templateData struct {
	Snippet models.Snippet
	Snippets []models.Snippet
	CurrentYear int
	Form any
	Flash string // Add a Flash field to the templateData struct
	IsAuthenticated bool // Add an IsAuthenticated field to the templateData struct
	CSRFToken string
}



