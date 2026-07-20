package middleware

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"
)

type AccessLog struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	RouteID   string `json:"route_id"`
	Status    int    `json:"status"`
	Bytes     int    `json:"bytes"`
	Duration  string `json:"duration"`
	Backend   string `json:"backend"`
	ClientIP  string `json:"client_ip"`
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		entry := AccessLog{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			RequestID: GetRequestID(r.Context()),
			Method:    r.Method,
			Path:      r.URL.Path,
			RouteID:   "-", // replace later
			Status:    rw.Status,
			Bytes:     rw.Bytes,
			Duration:  time.Since(start).String(),
			Backend:   "-", // replace after proxy
			ClientIP:  ip,
		}

		b, _ := json.Marshal(entry)

		log.Println(string(b))
	})
}
