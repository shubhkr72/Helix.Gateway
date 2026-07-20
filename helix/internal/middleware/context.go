package middleware

import "context"

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
)

func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDKey).(string)
	return id
}
