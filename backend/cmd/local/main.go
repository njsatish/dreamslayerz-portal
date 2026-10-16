package main

import (
	"github.com/njsatish/dreamslayerz-portal/backend/internal/api"
	appconfig "github.com/njsatish/dreamslayerz-portal/backend/internal/config"
	"log"
	"net/http"
)

func main() {
	cfg, err := appconfig.Load("config")
	if err != nil {
		log.Fatal(err)
	}
	handler := cors(api.New(cfg))
	log.Println("Dream Slayerz API: http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://localhost:5174" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
