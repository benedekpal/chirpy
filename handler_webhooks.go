package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/benedekpal/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	UserUpgradeEvent = "user.upgraded"
)

func (cfg *apiConfig) udpadeUserStatus(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != UserUpgradeEvent {
		respondWithError(w, http.StatusNoContent, "Invalid event", err)
		return
	}

	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't cast id to UUID", err)
		return
	}

	_, err = cfg.db.SetChirpyRedStatus(r.Context(), database.SetChirpyRedStatusParams{
		IsChirpyRed: true,
		ID:          id,
	})

	if errors.Is(err, pgx.ErrNoRows) { // or errors.Is(err, pgx.ErrNoRows)
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
