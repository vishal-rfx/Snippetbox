package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vishal-rfx/snippetbox/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the
// web application.
type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
	users *models.UserModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
}


func main(){
	// Define a new command line flag with the name 'addr', a default value of ':4000' and
	// some short help text explaining what the flag controls. The value of the flag will be stored
	// in the addr variable at the runtime
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Importantly, we use the flag.Parse() function to parse the command-line flag. 
	// This reads in the command-line flag value and assigns it to the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.

	// Define a new command line flag for the MySQL DSN string
	dsn := flag.String("dsn", "web:vishal@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	}))

	db, err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour


	app := &application{
		templateCache: templateCache,
		logger: logger,
		snippets: &models.SnippetModel{DB : db},
		users: &models.UserModel{DB: db},
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we want the server to use. In this case
	// the only thing that we're changing is the curve preferences value, so that only elliptic curves with assembly implementations 
	// are used
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before.
	srv := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		// Create a *log.Logger from our structured logger handler, which writes log entries at Error level,
		// and assign it to the ErrorLog field. If you would prefer to log the server errors at Warn level instead,
		// you could pass slog.LevelWarn as the final parameter.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("Starting a server on %s", "addr",*addr)
	
	
	// Use the ListenAndServeTLS method to start the HTTPS server. We pass in the paths to the TLS certificate and corresponding
	// private key as the two parameters.
	// Note that any error returned by ListenAndServe is always non nil
	// Each time the server receives a new HTTP request it will pass the request on to 
	// the servermux and in turn the servemux will check the URL path and dispatch the request
	// to the matching handler
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

// openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

