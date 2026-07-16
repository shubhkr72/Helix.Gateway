package handlers

import (
	"net/http"

	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/router"
)

type Gateway struct {
	Config *config.Config
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Gateway-owned endpoints
	switch r.URL.Path {
	case "/":
		Home(w, r)
		return

	case "/healthz":
		Healthz(w, r)
		return

	case "/readyz":
		Readyz(w, r)
		return

	case "/allservices":
		AllServices(w, r, g.Config)
	return
	}

	// Match configured routes
	match := router.MatchRoute(g.Config.Routes, r.URL.Path)

	if !match.Found {
		WriteError(
			w,
			http.StatusNotFound,
			map[string]any{
				"error": "route not found",
				"method":r.Method,
				"path":  r.URL.Path,
			},
		)
		return
	}

	// Day 4: Replace this with reverse proxy
	w.Write([]byte("Matched " + match.Route.ID + " -> " + match.Path))
}
