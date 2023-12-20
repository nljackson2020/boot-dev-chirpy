package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		Author_id: dbChirp.Author_id,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	authorID := -1
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = strconv.Atoi(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if authorID != -1 && dbChirp.Author_id != authorID {
			continue
		}
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			Body:      dbChirp.Body,
			Author_id: dbChirp.Author_id,
		})
	}

	sortChirp := r.URL.Query().Get("sort")
	if sortChirp == "asc" || sortChirp == "" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})

		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID > chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)

}
