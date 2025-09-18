package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/benedekpal/chirpy/internal/auth"
	"github.com/benedekpal/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (a *apiConfig) handlerAddUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email          string `json:"email"`
		HashedPassword string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Email) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "email is too long", nil)
		return
	}

	hashedPassWord, hErr := auth.HashPassword(params.HashedPassword)
	if hErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", hErr)
		return
	}

	u, dbErr := a.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hashedPassWord,
			Valid:  true,
		},
	})
	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new user", dbErr)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			Email:     u.Email,
		},
	})
}

func (a *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not provided", err)
		return
	}

	userid, err := auth.ValidateJWT(token, a.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not provided", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassWord, hErr := auth.HashPassword(params.Password)
	if hErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", hErr)
		return
	}
	if hashedPassWord == "" {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", hErr)
		return
	}

	updatedUser, err := a.db.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		HashedPassword: sql.NullString{
			String: hashedPassWord,
			Valid:  true,
		},
		Email: params.Email,
		ID:    userid,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email:     updatedUser.Email,
		},
	})
}
