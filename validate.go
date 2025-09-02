package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {

	type chirp struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
		Body  string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")

	respBody := returnVals{}

	decoder := json.NewDecoder(r.Body)
	ch := chirp{}
	err := decoder.Decode(&ch)

	if err != nil {
		respBody.Valid = false
		respBody.Error = "Something went wrong"
		respBody.Body = ""
		w.WriteHeader(500)
	} else if len(ch.Body) > 140 {
		respBody.Valid = false
		respBody.Error = "Chirp is too long"
		respBody.Body = ""
		w.WriteHeader(400)
	} else {
		respBody.Valid = true
		respBody.Error = ""
		respBody.Body = ""
		w.WriteHeader(200)
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Write(dat)

}
