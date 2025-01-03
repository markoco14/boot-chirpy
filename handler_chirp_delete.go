package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// get bearer token
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// validate token
	tokenUUID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID is required", nil)
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID format", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}

	if chirp == (database.Chirp{}) {
		respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
		return
	}

	if chirp.UserID != tokenUUID {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}
	
	deleteParams := database.DeleteChirpParams{
		ID:     chirpUUID,
		UserID: tokenUUID,
	}

	err = cfg.db.DeleteChirp(r.Context(), deleteParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}