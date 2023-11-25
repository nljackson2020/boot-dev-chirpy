package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginParams struct {
	Password           string `json:"password"`
	Email              string `json:"email"`
	Expires_in_seconds *int   `json:"expires_in_seconds,omitempty"`
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := LoginParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	db, err := NewDB("./database.json")
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side creating database")
		return
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side accessing database")
		return
	}

	verifiedUser := UserResponse{}

	for _, user := range dbStructure.User {
		if params.Email == user.Email {
			hashedPassword := user.Password
			plaintextPassword := params.Password

			err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))

			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Incorrect password")
				return
			}

			verifiedUser.Email = user.Email
			verifiedUser.ID = user.ID
		}
	}
	respondWithJSON(w, http.StatusOK, verifiedUser)
}
