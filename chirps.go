package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type parameters struct {
	Body string `json:"body"`
}

func handlerChirpsPostData(w http.ResponseWriter, r *http.Request) {

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

	db, err := NewDB("./database.json")
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side creating database")
	}

	dbErr := db.ensureDB()
	if dbErr != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Database not created")
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side accessing database")
		return
	}

	maxID := 0
	for id := range dbStructure.Chirps {
		if id > maxID {
			maxID = id
		}
	}

	newChirp := Chirp{
		Body: clean,
		ID:   maxID + 1,
	}

	dbStructure.Chirps[newChirp.ID] = newChirp
	err = db.writeDB(dbStructure)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error writing data to database")
		return
	}

	respondWithJSON(w, http.StatusCreated, newChirp)
}

func handlerChirpsGetData(w http.ResponseWriter, r *http.Request) {
	db, err := NewDB("./database.json")
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side creating database")
	}

	dbErr := db.ensureDB()
	if dbErr != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Database not created")
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
	}

	chirp, err := db.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	// var chirps []Chirp
	// chirps, err = db.GetChirps()
	// if err != nil {
	// 	respondWithError(w, http.StatusServiceUnavailable, "Can't access database")
	// }

	respondWithJSON(w, http.StatusOK, chirp)

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
