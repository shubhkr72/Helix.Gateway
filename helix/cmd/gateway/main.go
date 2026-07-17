package main

import (
	"log"
	"net/http"

	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/gateway"
	"github.com/shubhkr72/helix/internal/handlers"
	"github.com/shubhkr72/helix/internal/middleware"
)

func main() {
	cfg, err := config.Load("configs/gateway.yaml")
	if err != nil {
		log.Fatal(err)
	}

	gateway.PrintBanner(cfg)

	// Base handler
	handler := &handlers.Gateway{
		Config: cfg,
	}

	// Wrap with middleware
	handlerWithMiddleware := middleware.RequestID(handler)

	log.Fatal(http.ListenAndServe(
		":8080",
		handlerWithMiddleware,
	))
}