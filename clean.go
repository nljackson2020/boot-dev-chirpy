package main

import "strings"

func cleanChirp(s string) string {
	chirp := strings.Split(s, " ")
	badWords := [3]string{"kerfuffle", "sharbert", "fornax"}

	for i, word := range chirp {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				chirp[i] = "****"
			}
		}
	}

	return strings.Join(chirp, " ")
}
