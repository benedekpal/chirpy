package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "./public"

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create the http.Server struct
	server := &http.Server{
		Addr:    ":" + port, // Listen on port 8080
		Handler: mux,        // Use the custom ServeMux as the handler
	}

	// Serve static files
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	err := server.ListenAndServe()
	//log.Fatal(srv.ListenAndServe())

	if err != nil {
		log.Fatalf("error closing server: %v", err)
	}
}
