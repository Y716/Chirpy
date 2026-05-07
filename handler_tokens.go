package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Y716/chirpy/internal/auth"
	"github.com/Y716/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, "Unathorized", err)
		return
	}

	tokenData, err := apiCfg.database.GetUserFromRefreshToken(r.Context(), userToken)
	if err != nil || time.Now().After(tokenData.ExpiresAt) || tokenData.RevokedAt.Valid != false {
		fmt.Println(tokenData.RevokedAt.Valid != true)
		RespondWithError(w, 401, "Couldn't get token", err)
		return
	}

	token, err := auth.MakeJWT(tokenData.UserID, apiCfg.jwtSecret, time.Duration(1)*time.Hour)
	if err != nil {
		RespondWithError(w, 401, "Unable to make JWT Token", err)
		return
	}

	type parameters struct {
		Token string `json:"token"`
	}

	RespondWithJson(w, 200, parameters{Token: token})

}

func (apiCfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, "Unathorized", err)
		return
	}

	tokenData, err := apiCfg.database.GetUserFromRefreshToken(r.Context(), userToken)
	if err != nil || time.Now().After(tokenData.ExpiresAt) || tokenData.RevokedAt.Valid != false {
		RespondWithError(w, 401, "Couldn't get token", err)
		return
	}

	args := database.RevokeATokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: time.Now(),
		Token:     tokenData.Token,
	}

	apiCfg.database.RevokeAToken(r.Context(), args)

	RespondWithJson(w, 204, nil)

}
