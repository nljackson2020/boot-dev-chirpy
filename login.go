package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	db, err := NewDB("./database.json")
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side creating database")
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side accessing database")
		return
	}

	for _, user := range dbStructure.User {
		if params.Email == user.Email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))

			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Incorrect Password")
				return
			}

			returnUser := ReturnUser{
				ID:    user.ID,
				Email: user.Email,
			}
			respondWithJSON(w, http.StatusOK, returnUser)
		}
	}
}
