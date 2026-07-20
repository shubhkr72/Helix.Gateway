package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/shubhkr72/helix/internal/auth"
	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/jwt"
	"github.com/shubhkr72/helix/internal/router"
)

func Authentication(cfg *config.Config, manager *jwt.Manager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		match := router.MatchRoute(cfg.Routes, r.URL.Path)

		if !match.Found {
			next.ServeHTTP(w, r)
			return
		}

		if match.Route.Public {
			next.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := manager.VerifyToken(parts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		principal := &auth.Principal{
			UserID: claims.Subject,
			Email:  claims.Email,
			Roles:  claims.Roles,
		}

		if len(match.Route.Roles) > 0 {
			if !hasRole(principal, match.Route.Roles) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}

		ctx := auth.SetPrincipal(r.Context(), principal)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func hasRole(principal *auth.Principal, roles []string) bool {
	if principal == nil {
		return false
	}

	for _, required := range roles {
		for _, actual := range principal.Roles {
			if required == actual {
				return true
			}
		}
	}

	return false
}

func bearerToken(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)

	if len(parts) != 2 {
		return "", errors.New("invalid authorization header")
	}

	if parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	return parts[1], nil
}
