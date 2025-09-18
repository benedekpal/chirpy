package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/benedekpal/chirpy/internal/auth"
	"github.com/benedekpal/chirpy/internal/database"
)

const (
	DefaultExpiryJWT = 3600
	DefaultExpiryRT  = 60 * 24 * time.Hour
)

func (a *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := a.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user", err)
		return
	}

	errPassword := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword.String)
	if errPassword != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", errPassword)
		return
	}

	// Generate JWT
	jwtToken, err := auth.MakeJWT(dbUser.ID, a.jwtSecret, time.Duration(DefaultExpiryJWT)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create JWT", err)
		return
	}

	// Generate Refresh Token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create Refresh Token", err)
		return
	}
	refreshTokenFromDB, err := a.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(DefaultExpiryRT),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not save refreshToken", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
		Token:        jwtToken,
		RefreshToken: refreshTokenFromDB.Token,
	})
}
