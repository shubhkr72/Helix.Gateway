package gateway

import (
	"fmt"
	"github.com/shubhkr72/helix/internal/config"
)

func PrintBanner(cfg *config.Config) {

	fmt.Println("=========================================================")
	fmt.Println("                 HELIX API GATEWAY")
	fmt.Println("=========================================================")
	fmt.Println()

	fmt.Printf("Listening on :%d\n\n", cfg.Server.Port)

	fmt.Println("Gateway Endpoints")
	fmt.Println("-----------------")
	fmt.Println("GET    /")
	fmt.Println("GET    /healthz")
	fmt.Println("GET    /readyz")
	fmt.Println("GET    /allservices")
	fmt.Println()

	fmt.Println("Configured Routes")
	fmt.Println("---------------------------------------------------------")
	fmt.Printf("%-12s %-15s %-10s %-10s\n",
		"ID",
		"PATH",
		"STRIP",
		"BACKENDS",
	)

	fmt.Println("---------------------------------------------------------")

	for _, route := range cfg.Routes {

		fmt.Printf(
			"%-12s %-15s %-10t %-10d\n",
			route.ID,
			route.Path,
			route.StripPrefix,
			len(route.Backend),
		)
	}

	fmt.Println("---------------------------------------------------------")
	fmt.Println()
	fmt.Println("Gateway started successfully")
	fmt.Println("=========================================================")
}
