package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

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

		result, err := limiter.Allow(r.Context(), key)
		if err != nil {

			switch match.Route.RateLimit.FailurePolicy {

			case "fail_open":
				// Allow the request if Redis is unavailable.
				next.ServeHTTP(w, r)
				return

			case "fail_closed":
				fallthrough

			default:
				// Reject protected routes.
				writeError(
					w,
					http.StatusServiceUnavailable,
					"rate limiter unavailable",
				)
				return
			}
		}

		// Standard Rate Limit Headers
		w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(int64(result.ResetAfter.Seconds()), 10))

		if !result.Allowed {
			w.Header().Set(
				"Retry-After",
				strconv.FormatInt(int64(result.RetryAfter.Seconds()), 10),
			)

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
