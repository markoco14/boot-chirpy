package main

import (
	"chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerUsersRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	requestToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get token", err)
		return
	}

	// Now I have the refresh tooken, I need to use it to get the user.
	// I don't need to go to the refresh token table, I can query directly on the users table
	// no I don't need to get the user, I just need to get the refresh token, it has the UUID
	refreshToken, err := cfg.db.GetRefreshTokenByToken(r.Context(), requestToken)
	// handle case where token not found in DB
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	// handle case where token revoked
	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token revoked", nil)
		return
	}
	// handle case where token expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Token expired", nil)
		return
	}
	// token good to go, create new access token
	expiresIn := time.Hour
	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	responseJSON := response{
		Token: accessToken,
	}

	respondWithJSON(w, http.StatusOK, responseJSON)
}