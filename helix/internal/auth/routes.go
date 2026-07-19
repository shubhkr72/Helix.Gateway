package auth

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
) {
	mux.HandleFunc("/", handler.HealthCheck)
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/login", handler.Login)
}
