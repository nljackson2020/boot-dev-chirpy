package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	clean := cleanChirp(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: clean,
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

// func handlerValidate(w http.ResponseWriter, r *http.Request) {
// 	type parameters struct {
// 		Body string `json:"body"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		log.Printf("Error decoding parameters: %s", err)
// 		w.WriteHeader(500)
// 		return
// 	}

// 	bodyLength := len(params.Body)
// 	if bodyLength > 140 {
// 		type returnErr struct {
// 			Error string `json:"error"`
// 		}

// 		respBody := returnErr{
// 			Error: "Chirp is too long",
// 		}

// 		data, err := json.Marshal(respBody)
// 		if err != nil {
// 			log.Printf("Error marshalling JSON: %s", err)
// 			w.WriteHeader(500)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(400)
// 		w.Write(data)
// 		return
// 	}

// 	type returnVals struct {
// 		Valid bool `json:"valid"`
// 	}

// 	respBody := returnVals{
// 		Valid: true,
// 	}

// 	data, err := json.Marshal(respBody)
// 	if err != nil {
// 		log.Printf("Error marshalling JSON: %s", err)
// 		w.WriteHeader(500)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write(data)
// }
