package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/middleware"
)

func AllServices(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	w.Header().Set("Content-Type", "application/json")

	type RouteInfo struct {
		ID          string   `json:"id"`
		Path        string   `json:"path"`
		StripPrefix bool     `json:"strip_prefix"`
		Backends    []string `json:"backends"`
	}

	routes := make([]RouteInfo, 0, len(cfg.Routes))

	for _, route := range cfg.Routes {
		routes = append(routes, RouteInfo{
			ID:          route.ID,
			Path:        route.Path,
			StripPrefix: route.StripPrefix,
			Backends:    route.Backend,
		})
	}

	json.NewEncoder(w).Encode(map[string]any{
		"service": "Helix Gateway",
		"routes":  routes,
	})
}
func Home(w http.ResponseWriter, r *http.Request) {
    // log.Println("Request ID:", id)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"service": "helix gateway",
		"status":  "running",
	})
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	id := middleware.GetRequestID(r.Context())

	// log.Println("Request ID:", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status":     "ok",
		"request_id": id,
	})
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	id := middleware.GetRequestID(r.Context())

	// log.Println("Request ID:", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"status":     "ready",
		"request_id": id,
	})
}
