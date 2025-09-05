package main

import (
	"fmt"
	"net/http"
	"os"
)

func (cfg *apiConfig) readMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data, err := os.ReadFile("./local/metrics_template.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Error reading file:", err)
		return
	}

	result := fmt.Sprintf(string(data), cfg.fileserverHits.Load())

	_, err = w.Write([]byte(result))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	err := cfg.db.ClearUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Couldn't clear users table", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("fileserverHits reseted to 0\nCleared Users table"))
}
