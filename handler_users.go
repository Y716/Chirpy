package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Y716/chirpy/internal/auth"
	"github.com/Y716/chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	args := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := apiCfg.database.CreateUser(r.Context(), args)
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

func (apiCfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnBody struct {
		Id            uuid.UUID `json:"id"`
		Created_at    time.Time `json:"created_at"`
		Updated_at    time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Token         string    `json:"token"`
		Refresh_Token string    `json:"refresh_token"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := apiCfg.database.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	flag, err := auth.ComparePasswordHash(params.Password, user.HashedPassword)
	if err != nil || !flag {
		RespondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, apiCfg.jwtSecret, time.Duration(1)*time.Hour)
	if err != nil {
		RespondWithError(w, 401, "Unable to make JWT Token", err)
		return
	}

	refresh_token := auth.MakeRefreshToken()

	args := database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{},
	}

	refreshTokenData, err := apiCfg.database.CreateRefreshToken(r.Context(), args)
	if err != nil {
		RespondWithError(w, 401, "Unable to make Refresh Token", err)
		return
	}

	returnVal := returnBody{
		Id:            user.ID,
		Created_at:    user.CreatedAt,
		Updated_at:    user.UpdatedAt,
		Email:         user.Email,
		Token:         token,
		Refresh_Token: refreshTokenData.Token,
	}

	RespondWithJson(w, 200, returnVal)

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
