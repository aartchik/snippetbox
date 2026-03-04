package main

import (
	"database/sql"
	"flag"
	"html/template"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.net/internal/models"
	"github.com/go-playground/form/v4"
	"crypto/tls"
	"log"
	"net/http"
	"os"

	 "github.com/alexedwards/scs/mysqlstore" 
    "github.com/alexedwards/scs/v2"    
	"time" 
)

type config struct {
	addr      string
	staticDir string
	dsn       string
	debug     bool
	tls	      bool
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets  models.SnippetModelInterface
	users models.UserModelInterface
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	sessionManager *scs.SessionManager
	debug bool
}

func openDB(cfg *config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "Database connection string")
	flag.BoolVar(&cfg.debug, "debug", false, "When running in debug mode, any detailed errors and stack traces should be displayed in the browser")
	flag.BoolVar(&cfg.tls, "tls", false, "Enable HTTPS")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(&cfg)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 24 * time.Hour * 7
	sessionManager.Cookie.Secure = cfg.tls

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
		users: &models.UserModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
		debug: cfg.debug,
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		TLSConfig: tlsConfig,
	}

	infoLog.Printf("Starting server on %s", cfg.addr)
	if cfg.tls {
		err = srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem")
	} else {
		err = srv.ListenAndServe()
	}
	errorLog.Fatal(err)
}

