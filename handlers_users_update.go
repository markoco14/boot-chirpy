package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)


func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email	string `json:"email"`
		Password 		string `json:"password"`
	}

	type response struct {
		User
	}

	// get access token from Authorization header
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No access token", err)
		return
	}
	// validate access token and return user UUID
	userUUID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}	
	// decode request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	// validate the email
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required", nil)
		return
	}

	// validate the password
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required", nil)
		return
	}
	if len(params.Password) <= 5 {
		respondWithError(w, http.StatusBadRequest, "Password must be at least 6 characters", nil)
		return
	}

	HashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "An error occurred", err)
		return
	}

	UpdateParams := database.UpdateUserEmailPasswordParams{
		Email: params.Email,
		HashedPassword: HashedPassword,
		ID: userUUID,
	}

	UpdatedUser, err := cfg.db.UpdateUserEmailPassword(r.Context(), UpdateParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	responseUser := User{
		ID:        UpdatedUser.ID,
		CreatedAt: UpdatedUser.CreatedAt,
		UpdatedAt: UpdatedUser.UpdatedAt,
		Email:     UpdatedUser.Email,
		IsChirpyRed: UpdatedUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, response{User: responseUser})
}