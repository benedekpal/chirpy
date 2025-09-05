package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/benedekpal/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	ready          atomic.Bool
	dbURL          string
	platform       string
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

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	dbQueries := database.New(db)

	const port = "8080"
	const filepathRoot = "./public"

	//app := &App{}

	config := &apiConfig{}
	config.fileserverHits.Store(0)
	config.db = dbQueries
	config.platform = platform
	config.dbURL = dbURL

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create the http.Server struct
	server := &http.Server{
		Addr:    ":" + port, // Listen on port 8080
		Handler: mux,        // Use the custom ServeMux as the handler
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", config.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", config.handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", config.handlerChirps)
	mux.HandleFunc("POST /api/users", config.handlerAddUser)
	mux.HandleFunc("GET /admin/metrics", config.readMetrics)
	mux.HandleFunc("POST /admin/reset", config.resetMetrics)
	mux.HandleFunc("POST /api/chirps", config.validateAndSaveChirp)
	mux.HandleFunc("GET /api/chirps", config.retrieveAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", config.retrieveChirpByID)

	// after init:
	//app.ready.Store(true)
	config.ready.Store(true)

	err = server.ListenAndServe()
	//log.Fatal(srv.ListenAndServe())

	if err != nil {
		log.Fatalf("error closing server: %v", err)
	}
}
