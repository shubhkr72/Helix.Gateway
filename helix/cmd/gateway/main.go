package main

import (
	"log"
	"net/http"

	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/gateway"
	"github.com/shubhkr72/helix/internal/handlers"
	"github.com/shubhkr72/helix/internal/jwt"
	"github.com/shubhkr72/helix/internal/middleware"
	"github.com/shubhkr72/helix/internal/proxy"
)

func main() {
	cfg, err := config.Load("configs/gateway.yaml")
	if err != nil {
		log.Fatal(err)
	}

	gateway.PrintBanner(cfg)

	jwtManager, err := jwt.NewGatewayManager(
		cfg.JWT.PublicKey,
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
	)
	if err != nil {
		log.Fatal(err)
	}

	proxies := make(map[string]*proxy.Gateway)

	for _, route := range cfg.Routes {
		if len(route.Backend) == 0 {
			log.Fatalf("route %q has no backend configured", route.ID)
		}

		p, err := proxy.New(route.Backend[0])
		if err != nil {
			log.Fatalf("route %q: %v", route.ID, err)
		}

		proxies[route.ID] = p
	}

	handler := &handlers.Gateway{
		Config:  cfg,
		Proxies: proxies,
	}

	handlerWithMiddleware := middleware.RequestID(
		middleware.Authentication(
			cfg,
			jwtManager,
			middleware.Logging(
				handler,
			),
		),
	)

	log.Println("Gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddleware))
}
