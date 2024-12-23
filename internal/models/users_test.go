package models

import (
	"testing"

	"github.com/vishal-rfx/snippetbox/internal/assert"
)

func TestUserModelExits(t *testing.T) {
    // Setup a suite of table-driven tests and expected results.
    tests := []struct {
        name string
        userID int
        want bool
    }{
        {
            name: "Valid ID",
            userID: 1,
            want: true,
        },
        {
            name: "Zero ID",
            userID: 0,
            want: false,
        },
        {
            name: "Non-existent ID",
            userID: 2,
            want: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T){
            // Call the newTestDB() helper function to get a connection pool to our database. Calling this here inside t.Run()
            // means that fresh database tables and data will be set up and torn down for each sub test
            db := newTestDB(t)
            // Create a new instance of the UserModel
            m := UserModel{db}

            exists, err := m.Exists(tt.userID)
            assert.Equal(t, exists, tt.want)
            assert.NilError(t, err)
        })
    }
}