package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "You can't access that reset", nil)
		return
	}
	
	err := cfg.db.DeleteUsers(r.Context())
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete users", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
