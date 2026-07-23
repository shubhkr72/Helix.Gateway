package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	Timeouts struct {
		Dial           time.Duration `yaml:"dial"`
		ResponseHeader time.Duration `yaml:"response_header"`
	} `yaml:"timeouts"`

	JWT JWTConfig `yaml:"jwt"`

	Redis RedisConfig `yaml:"redis"`

	Routes []Route `yaml:"routes"`
}

type JWTConfig struct {
	Issuer    string `yaml:"issuer"`
	Audience  string `yaml:"audience"`
	PublicKey string `yaml:"public_key"`
}

type RateLimitConfig struct {
	Enabled       bool    `yaml:"enabled"`
	Capacity      float64 `yaml:"capacity"`
	RefillRate    float64 `yaml:"refill_rate"`
	KeyStrategy   string  `yaml:"key_strategy"`
	FailurePolicy string  `yaml:"failure_policy"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
}

type Route struct {
	ID            string          `yaml:"id"`
	Path          string          `yaml:"path"`
	Backend       []string        `yaml:"backend"`
	StripPrefix   bool            `yaml:"strip_prefix"`
	Public        bool            `yaml:"public"`
	Authenticated bool            `yaml:"authenticated"`
	Roles         []string        `yaml:"roles"`
	RateLimit     RateLimitConfig `yaml:"rate_limit"`
}

func Load(path string) (*Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := Config{}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	if cfg.Timeouts.Dial == 0 {
		cfg.Timeouts.Dial = 5 * time.Second
	}

	if cfg.Timeouts.ResponseHeader == 0 {
		cfg.Timeouts.ResponseHeader = 10 * time.Second
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validate(cfg Config) error {

	ids := map[string]bool{}

	for _, r := range cfg.Routes {

		if r.ID == "" {
			return errors.New("route id required")
		}

		if ids[r.ID] {
			return fmt.Errorf("duplicate route id %s", r.ID)
		}
		ids[r.ID] = true

		if len(r.Backend) == 0 {
			return fmt.Errorf("%s has no backend", r.ID)
		}

		for _, b := range r.Backend {
			u, err := url.Parse(b)
			if err != nil || u.Scheme == "" || u.Host == "" {
				return fmt.Errorf("invalid backend %s", b)
			}
		}

		if r.RateLimit.Capacity <= 0 {
			return fmt.Errorf("%s has invalid rate limit capacity", r.ID)
		}

		if r.RateLimit.RefillRate <= 0 {
			return fmt.Errorf("%s has invalid refill rate", r.ID)
		}

		switch r.RateLimit.KeyStrategy {
		case "user", "ip", "api_key", "global":
		default:
			return fmt.Errorf("%s has invalid key strategy", r.ID)
		}
	}

	paths := []string{}

	for _, r := range cfg.Routes {
		if slices.Contains(paths, r.Path) {
			return fmt.Errorf("duplicate path %s", r.Path)
		}
		paths = append(paths, r.Path)
	}

	return nil
}
