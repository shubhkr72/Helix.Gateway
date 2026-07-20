package middleware

import "context"

type contextKey string

const RequestIDKey contextKey = "request_id"

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
