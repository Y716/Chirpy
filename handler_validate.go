package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		RespondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	cleanedString := censorProfanity(params.Body)

	RespondWithJson(w, 200, returnBody{
		CleanedBody: cleanedString,
	})
}

func censorProfanity(s string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}

	splittedString := strings.Split(s, " ")
	for i, word := range splittedString {
		lowered := strings.ToLower(word)
		if ok := slices.Contains(profanity, lowered); ok {
			splittedString[i] = "****"
		}
	}

	return strings.Join(splittedString, " ")
}
