package main

import (
	"encoding/json"
	"net/http"

	"github.com/nljackson2020/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserID int `json:"user_id"`
	}

	type PolkaWebhook struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := PolkaWebhook{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode webhook")
		return
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find ApiKey")
		return
	}

	if apiKey != cfg.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find ApiKey")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, "")
		return
	}

	err = cfg.DB.UpdateChirpyRed(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	respondWithJSON(w, http.StatusOK, "")
}
