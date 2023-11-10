package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

var cfg config

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	flag.StringVar(&cfg.addr, "addr", ":4000", "Http network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	app := &application{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	app.infoLog.Printf("Starting server on port: %s\n", cfg.addr)
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLog,
		Handler:  mux,
	}
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
