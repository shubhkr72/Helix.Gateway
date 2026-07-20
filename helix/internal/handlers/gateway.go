package handlers

import (
	"net/http"

	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/jwt"
	"github.com/shubhkr72/helix/internal/proxy"
	"github.com/shubhkr72/helix/internal/router"
)

type Gateway struct {
	Config  *config.Config
	JWT     *jwt.Manager
	Proxies map[string]*proxy.Gateway
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	match := router.MatchRoute(g.Config.Routes, r.URL.Path)

	if !match.Found {
		WriteError(
			w,
			http.StatusNotFound,
			map[string]any{
				"error":  "route not found",
				"method": r.Method,
				"path":   r.URL.Path,
			},
		)
		return
	}

	p, ok := g.Proxies[match.Route.ID]
	if !ok {
		WriteError(
			w,
			http.StatusServiceUnavailable,
			map[string]any{
				"error": "backend unavailable",
			},
		)
		return
	}

	r.URL.Path = match.Path

	p.ServeHTTP(w, r)
}
