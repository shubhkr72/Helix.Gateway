package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/ratelimiter"
	"github.com/shubhkr72/helix/internal/router"
)

func RateLimit(
	cfg *config.Config,
	limiters map[string]ratelimiter.Limiter,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		match := router.MatchRoute(cfg.Routes, r.URL.Path)

		if !match.Found {
			next.ServeHTTP(w, r)
			return
		}

		limiter, ok := limiters[match.Route.ID]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		var userID string

		if v := r.Context().Value("userID"); v != nil {
			if s, ok := v.(string); ok {
				userID = s
			}
		}

		key := ratelimiter.BuildKey(
			r,
			match.Route.RateLimit.KeyStrategy,
			userID,
		)

		allowed, err := limiter.Allow(r.Context(), key)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "rate limiter failure")
			return
		}

		if !allowed {
			writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": msg,
	})
}
