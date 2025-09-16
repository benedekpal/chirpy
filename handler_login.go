package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/benedekpal/chirpy/internal/auth"
)

const (
	DefaultExpirySeconds = 3600
)

func (a *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expiry   int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Expiry == 0 || params.Expiry > DefaultExpirySeconds {
		params.Expiry = DefaultExpirySeconds
	}

	dbUser, err := a.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user", err)
		return
	}

	errPassword := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword.String)
	if errPassword != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", errPassword)
	}

	// Generate JWT using expires_in_seconds
	token, err := auth.MakeJWT(dbUser.ID, a.jwtSecret, time.Duration(params.Expiry)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
		Token: token,
	})
}
