package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"encoding/json"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	
	mux.HandleFunc("/api/validate_chirp", func(w http.ResponseWriter, r *http.Request){
		type parameters struct {
			Body string `json:"body"`
		}

		type responseBody struct {
			Valid bool `json:"valid"`
		}

		w.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Something went wrong"}`))
			return
		}
		
		if len(params.Body) > 140 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Chirp is too long"}`))
			return
		}
		resp := responseBody{Valid: true}
		data, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Something went wrong"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

