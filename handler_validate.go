package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

type Mask struct {
	UnAllowed   []string //:= ["kerfuffle", "sharbert", "fornax"]
	Replacement string
}

func NewProfaneWordsMask() Mask {
	return Mask{
		UnAllowed:   []string{"kerfuffle", "sharbert", "fornax"},
		Replacement: "****",
	}
}

func (m *Mask) censureString(s string) string {
	parts := strings.Split(s, " ")
	for i, word := range parts {
		if slices.Contains(m.UnAllowed, strings.ToLower(word)) {
			parts[i] = m.Replacement
		}
	}
	return strings.Join(parts, " ")
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body    string `json:"body"`
		User_id string `json:"user_id"`
	}
	type returnVals struct {
		Valid        bool   `json:"valid"`
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	mask := NewProfaneWordsMask()

	respondWithJSON(w, http.StatusOK, returnVals{
		Valid:        true,
		Cleaned_body: mask.censureString(params.Body),
	})

}
