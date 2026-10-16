package api

import (
	"encoding/json"
	appconfig "github.com/njsatish/dreamslayerz-portal/backend/internal/config"
	"net/http"
	"strings"
)

type Handler struct{ Config appconfig.Config }

func New(cfg appconfig.Config) Handler { return Handler{Config: cfg} }
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimSuffix(r.URL.Path, "/")
	switch {
	case r.Method == http.MethodGet && path == "/api/health":
		writeJSON(w, 200, map[string]string{"status": "ok", "service": "booksy-portal-api"})
	case r.Method == http.MethodGet && path == "/api/config":
		writeJSON(w, 200, h.Config)
	case r.Method == http.MethodGet && path == "/api/business":
		writeJSON(w, 200, h.Config.Business)
	case r.Method == http.MethodGet && path == "/api/providers":
		writeJSON(w, 200, h.Config.Providers)
	case r.Method == http.MethodGet && path == "/api/services":
		writeJSON(w, 200, h.Config.Services)
	case r.Method == http.MethodGet && path == "/api/availability":
		h.serveAvailability(w, r)
	default:
		writeJSON(w, 404, map[string]string{"error": "route not found"})
	}
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
