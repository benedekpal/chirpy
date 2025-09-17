package main

import (
	"net/http"
	"time"

	"github.com/benedekpal/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshAccessToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not provided", err)
		return
	}

	refreshTokenFromDB, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found", err)
		return
	}

	if refreshTokenFromDB.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token revoked", err)
		return
	}

	if refreshTokenFromDB.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Token expired", err)
		return
	}

	//generate new JWT
	jwtToken, err := auth.MakeJWT(refreshTokenFromDB.UserID, cfg.jwtSecret, time.Duration(DefaultExpiryJWT)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: jwtToken,
	})

}
