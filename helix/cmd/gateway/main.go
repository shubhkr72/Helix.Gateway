package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/gateway"
	"github.com/shubhkr72/helix/internal/handlers"
	"github.com/shubhkr72/helix/internal/jwt"
	"github.com/shubhkr72/helix/internal/middleware"
	"github.com/shubhkr72/helix/internal/proxy"
	"github.com/shubhkr72/helix/internal/ratelimiter"
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

	redisClient, err := ratelimiter.NewRedisClient(cfg.Redis.Addr)
	if err != nil {
		log.Printf("WARNING: Redis unavailable: %v", err)
	} else {
		defer redisClient.Close()
	}

	proxies := make(map[string]*proxy.Gateway)
	limiters := make(map[string]ratelimiter.Limiter)

	for _, route := range cfg.Routes {
		if len(route.Backend) == 0 {
			log.Fatalf("route %q has no backend configured", route.ID)
		}

		p, err := proxy.New(route.Backend[0])
		if err != nil {
			log.Fatalf("route %q: %v", route.ID, err)
		}

		proxies[route.ID] = p

		limiters[route.ID] = ratelimiter.NewRedisLimiter(
			redisClient,
			ratelimiter.Config{
				Capacity:    route.RateLimit.Capacity,
				RefillRate:  route.RateLimit.RefillRate,
				KeyStrategy: route.RateLimit.KeyStrategy,
			},
		)
	}

	handler := &handlers.Gateway{
		Config:  cfg,
		Proxies: proxies,
	}

	handlerWithMiddleware := middleware.RequestID(
		middleware.Authentication(
			cfg,
			jwtManager,
			middleware.RateLimit(
				cfg,
				limiters,
				middleware.Logging(handler),
			),
		),
	)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	log.Printf("Gateway listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, handlerWithMiddleware))
}
