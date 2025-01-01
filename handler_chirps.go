package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"chirpy/internal/database"
	"chirpy/internal/auth"

	"github.com/google/uuid"
)

type createChirpRequest struct {
	Body   string `json:"body"`
}

type Chirp struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
}

func dbChirpToResponse(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID.String(),
		CreatedAt: dbChirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: dbChirp.UpdatedAt.Format(time.RFC3339),
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID.String(),
	}
}

func (cfg *apiConfig) handleListChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	responseChirps := make([]Chirp, len(chirps))

	for i, chirp := range chirps {
		responseChirps[i] = dbChirpToResponse(chirp)
	}

	respondWithJSON(w, http.StatusOK, responseChirps)
}

func (cfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	parsedID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), parsedID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
		return
	}

	responseChirp := dbChirpToResponse(chirp)

	respondWithJSON(w, http.StatusOK, responseChirp)
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	tokenUUID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// code goes here
	decoder := json.NewDecoder(r.Body)
	params := createChirpRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid request payload", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := getCleanedBody(params.Body)


	// Create database params
	dbParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: tokenUUID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	responseChirp := dbChirpToResponse(chirp)

	respondWithJSON(w, http.StatusCreated, responseChirp)
}

func getCleanedBody(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

// func validateChirp(w http.ResponseWriter, r *http.Request) {
// 	type parameters struct {
// 		Body string `json:"body"`
// 	}
// 	type returnVals struct {
// 		CleanedBody string `json:"cleaned_body"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
// 		return
// 	}

// 	const maxChirpLength = 140
// 	if len(params.Body) > maxChirpLength {
// 		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
// 		return
// 	}

// 	cleanedBody := getCleanedBody(params.Body)

// 	respondWithJSON(w, http.StatusOK, returnVals{
// 		CleanedBody: cleanedBody,
// 	})
// }
