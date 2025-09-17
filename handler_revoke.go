package main

import (
	"net/http"

	"github.com/benedekpal/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not provided", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not revoke token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
