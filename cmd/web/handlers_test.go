package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/vishal-rfx/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}

func TestSnippetView(t *testing.T){
	// Create a new instance of the application struct which uses the mocked dependencies
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our application for different URLs
	tests := []struct {
		name string
		urlPath string
		wantCode int
		wantBody string
	}{
		{
			name: "Valid ID",
			urlPath: "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name: "Non-existent ID",
			urlPath: "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Negative ID",
			urlPath: "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Decimal ID",
			urlPath: "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name: "String ID",
			urlPath: "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Empty ID",
			urlPath: "/snippet/view",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}

}

func TestUserSignup(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set up the test
	// server for running an end-to-end test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET /user/signup request and then extract the CSRF token from the response body
	_, _, body := ts.get(t, "/user/signup")
	validCsrfToken := extractCSRFToken(t, body)

	// Set up some table-driven tests to check the responses sent by our application for different URLs
	const (
		validName = "Bob"
		validPassword = "validPassword"
		validEmail = "bob@example.com"
		formTag = "<form action='/user/signup' method='POST' novalidate>"
 	)

	tests := []struct {
		name string
		userName string
		userEmail string
		userPassword string
		csrfToken string
		wantCode int
		wantFormTag string
	}{
		{
			name : "Valid submission",
			userName: validName,
			userEmail: validEmail,
			userPassword: validPassword,
			csrfToken: validCsrfToken,
			wantCode: http.StatusSeeOther,
		},
		{
			name : "Invalid CSRF Token",
			userName: validName,
			userEmail: validEmail,
			userPassword: validPassword,
			csrfToken: "wrong-token",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Empty Name",
			userName: "",
			userEmail :validEmail,
			userPassword: validPassword,
			csrfToken: validCsrfToken,
			wantCode: http.StatusUnprocessableEntity,
			// wantFormTag: formTag,
		},
		{
			name: "Empty Email",
			userName: validName,
			userEmail: "",
			userPassword: validPassword,
			csrfToken: validCsrfToken,
			wantCode: http.StatusUnprocessableEntity,
			// wantFormTag: formTag,
		},
		{
			name: "Empty Password",
			userName: validName,
			userEmail: validEmail,
			userPassword: "",
			csrfToken: validCsrfToken,
			wantCode: http.StatusUnprocessableEntity,
			// wantFormTag: formTag,
		},
		{
			name: "Invalid Email",
			userName: validName,
			userEmail: "invalid-email",
			userPassword: validPassword,
			csrfToken: validCsrfToken,
			wantCode: http.StatusUnprocessableEntity,
			// wantFormTag: formTag,
		},
		{
			name: "Short password",
			userName : validName,
			userEmail: validEmail,	
			userPassword: "short",
			csrfToken: validCsrfToken,
			wantCode: http.StatusUnprocessableEntity,
			// wantFormTag: formTag,
		},
		{
			name: "Duplicate Email",
			userName: validName,
			userEmail: "dupe@example.com",
			userPassword: validPassword,
			csrfToken: validCsrfToken,
			wantCode: http.StatusUnprocessableEntity,
			// wantFormTag: formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}