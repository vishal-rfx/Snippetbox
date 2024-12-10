package main

import (
	"path/filepath"
	"time"
	"html/template"

	"github.com/vishal-rfx/snippetbox/internal/models"
)

// humanDate function which returns a nicely formatter string representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
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

	// use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl"
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Extract the filename from the full filepath and assign it to the name variable
		name := filepath.Base(page)
		// The template.FuncMap must be registered with the template set before we call the ParseFiles() method. This
		// means we have to use the template.New() to create an empty template set, use the Funcs method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() on this template set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() on this template set to  add the page template.
		ts, err = ts.ParseFiles(page)
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



