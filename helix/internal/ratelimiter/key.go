package ratelimiter

import (
	"net"
	"net/http"
	"strings"
)

const (
	KeyStrategyUser   = "user"
	KeyStrategyIP     = "ip"
	KeyStrategyAPIKey = "api_key"
	KeyStrategyGlobal = "global"
)

func BuildKey(r *http.Request, strategy string, userID string) string {
	switch strategy {
	case KeyStrategyUser:
		if userID == "" {
			return "anonymous"
		}
		return "user:" + userID

	case KeyStrategyAPIKey:
		apiKey := strings.TrimSpace(r.Header.Get("X-API-Key"))
		if apiKey == "" {
			return "anonymous"
		}
		return "apikey:" + apiKey

	case KeyStrategyIP:
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		return "ip:" + host

	case KeyStrategyGlobal:
		fallthrough
	default:
		return "global"
	}
}
