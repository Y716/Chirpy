package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type returnBody struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := apiCfg.database.CreateUser(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	returnVal := returnBody{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}

	RespondWithJson(w, 201, returnVal)
}

func (apiCfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if apiCfg.environment != "dev" {
		RespondWithError(w, 403, "Forbidden", nil)
		return
	}

	err := apiCfg.database.DeleteAllUsers(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't delete all users", err)
		return
	}
}
