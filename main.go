package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type App struct {
	ready atomic.Bool
}

func main() {
	const port = "8080"
	const filepathRoot = "./public"

	app := &App{}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create the http.Server struct
	server := &http.Server{
		Addr:    ":" + port, // Listen on port 8080
		Handler: mux,        // Use the custom ServeMux as the handler
	}

	// Serve static files
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", app.handlerReadiness)

	// after init:
	app.ready.Store(true)

	err := server.ListenAndServe()
	//log.Fatal(srv.ListenAndServe())

	if err != nil {
		log.Fatalf("error closing server: %v", err)
	}
}
