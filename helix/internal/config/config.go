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

	Routes []Route `yaml:"routes"`
}

type Route struct {
	ID          string   `yaml:"id"`
	Path        string   `yaml:"path"`
	Backend     []string `yaml:"backend"`
	StripPrefix bool     `yaml:"strip_prefix"`
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
