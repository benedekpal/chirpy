package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/benedekpal/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) validateAndSaveChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}
	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	userId, err := uuid.Parse(params.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode UserId", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	mask := NewProfaneWordsMask()

	newChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   mask.censureString(params.Body),
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save new chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      newChirp.Body,
			UserID:    newChirp.UserID,
		},
	})

}

func (cfg *apiConfig) retrieveAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirps := make([]Chirp, 0, len(dbChirps))
	for _, c := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: c.ID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
			Body: c.Body, UserID: c.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
