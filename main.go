package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create the http.Server struct
	server := &http.Server{
		Addr:    ":8080", // Listen on port 8080
		Handler: mux,     // Use the custom ServeMux as the handler
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
