package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type App struct {
	ready atomic.Bool
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	// little fun for showing uptime for myselfe
	defer func(start time.Time) {
		log.Printf("took %s", time.Since(start))
	}(time.Now())
	//equvivalent
	//start := time.Now()
	//defer func() {
	//	log.Printf("took %s", time.Since(start))
	//}()

	const port = "8080"
	const filepathRoot = "./public"

	app := &App{}

	config := &apiConfig{}
	config.fileserverHits.Store(0)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create the http.Server struct
	server := &http.Server{
		Addr:    ":" + port, // Listen on port 8080
		Handler: mux,        // Use the custom ServeMux as the handler
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", config.middlewareMetricsInc(handler))
	mux.HandleFunc("/healthz", app.handlerReadiness)
	mux.HandleFunc("/metrics", config.readMetrics)
	mux.HandleFunc("/reset", config.resetMetrics)

	// after init:
	app.ready.Store(true)

	err := server.ListenAndServe()
	//log.Fatal(srv.ListenAndServe())

	if err != nil {
		log.Fatalf("error closing server: %v", err)
	}
}
