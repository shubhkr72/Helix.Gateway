package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"regexp"
)

const HeaderRequestID = "X-Request-ID"

// allow only letters, numbers, dash and underscore
var validID = regexp.MustCompile(`^[A-Za-z0-9_-]{8,128}$`)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := r.Header.Get(HeaderRequestID)

		// Validate
		if !validID.MatchString(id) {
			id = newRequestID()
		}

		// Put into context
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		r = r.WithContext(ctx)

		// Echo back
		w.Header().Set(HeaderRequestID, id)

		next.ServeHTTP(w, r)
	})
}

func newRequestID() string {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return "unknown-request-id"
	}

	return hex.EncodeToString(b)
}