package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet type holds the data for an individual snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY) )
	`

	// Use the Exec() method on the embedded connection pool to execute the
	// statement, followed by the values for the placeholder parameters: title, content and expiry in that order.
	// This method returns a sql.Result type which contains some
	// basic information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result to get the ID of our newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
			 FROM snippets
			 WHERE expires > UTC_TIMESTAMP() and id = ?`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet
	// Use row.Scan() to copy the values from each field in sql.Row to the corresponding field in the Snippet struct.
	// The arguments to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of columns returned by your statement
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// Latest will return the slice of 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() 
		ORDER BY id DESC
		LIMIT 10
	`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// We defer rows.Close() to ensure the sql.Rows resultset is always properly closed before the Latest()
	// method returns. This defer statement should come after you check for an error from the Query() method,
	// Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	var snippets []Snippet
	// Use rows.Next to iterate through the rows in the resultset. This prepares the first (and then each subsequent)
	// row to be acted on by the rows.Scan() method. If iteration over all the rows completes then the resultset automatically 
	// closes itself and frees-up the underlying database connection.

	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any error
	// that was encountered during the iteration. It's important to call this - don't assume that a 
	// successful iteration was completed over the whole resultset.

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
