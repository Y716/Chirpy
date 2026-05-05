package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Y716/chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type returnBody struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		UserID     uuid.UUID `json:"user_id"`
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

	args := database.CreateChirpParams{
		Body:   cleanedString,
		UserID: params.UserID,
	}

	chirp, err := apiCfg.database.CreateChirp(r.Context(), args)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	returnVal := returnBody{
		Id:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		UserID:     chirp.UserID,
	}

	RespondWithJson(w, 201, returnVal)

}

func (apiCfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	type returnBody struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		UserID     uuid.UUID `json:"user_id"`
	}

	chirps, err := apiCfg.database.GetAllChirps(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	returnVal := []returnBody{}
	for _, chirp := range chirps {
		returnVal = append(returnVal, returnBody{
			Id:         chirp.ID,
			Created_at: chirp.CreatedAt,
			Updated_at: chirp.UpdatedAt,
			Body:       chirp.Body,
			UserID:     chirp.UserID,
		})
	}

	RespondWithJson(w, 200, returnVal)
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
