package main

import (
	"chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerUsersRevoke(w http.ResponseWriter, r *http.Request) {
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get token", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshTokenByToken(r.Context(), headerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// check if already revoked, return success
	if refreshToken.RevokedAt.Valid {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// revoke the token regardless
	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}

	// check if the token was expired anyways
	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Token expired", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
