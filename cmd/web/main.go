package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/notkisi/snippetbox/internal/models"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	config        *config
	snippets      *models.SnippetModel
	templateCache templCache
}

func main() {

	cfg := &config{}
	flag.StringVar(&cfg.addr, "addr", ":4000", "Http network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache := templCache{}
	templateCache.update()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		config:        cfg,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	go func() {
		// last modified = 1970
		lastModified := time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC)
		for true {
			fileInfo, _ := os.Stat("./ui/html/pages/home.tmpl")
			if fileInfo.ModTime().After(lastModified) {
				// todo for all templates
				app.infoLog.Println("Reloading templates")
				app.templateCache.update()
				lastModified = fileInfo.ModTime()
			}
			time.Sleep(time.Second)
		}
	}()

	app.infoLog.Printf("Starting server on port: %s\n", cfg.addr)
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
