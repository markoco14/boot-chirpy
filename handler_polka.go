package main

import (
	"chirpy/internal/database"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
)

type PolkaData struct {
	UserID string `json:"user_id"`
}

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data PolkaData `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	// check if the event is user.upgraded
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// parse the user_id
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user_id format", err)
		return
	}

    dbUser, err := cfg.db.GetUserByID(r.Context(), userID)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Couldn't get user", err)
        return
    }

    updateParams := database.UpdateUserRedMembershipParams{
        IsChirpyRed: true,
        ID: dbUser.ID,
    }

    _, err = cfg.db.UpdateUserRedMembership(r.Context(), updateParams)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}