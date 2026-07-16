package router

import (
	"testing"

	"github.com/shubhkr72/helix/internal/config"
)

func TestMatchRoute(t *testing.T) {

	routes := []config.Route{
		{
			ID:          "users",
			Path:        "/users",
			StripPrefix: true,
		},
		{
			ID:   "orders",
			Path: "/orders",
		},
	}

	tests := []struct {
		path string
		want bool
		out  string
	}{
		{"/users", true, "/"},
		{"/users/", true, "/"},
		{"/users/42", true, "/42"},
		{"/users-old", false, ""},
		{"/orders", true, "/orders"},
		{"/abc", false, ""},
	}

	for _, tt := range tests {

		m := MatchRoute(routes, tt.path)

		if m.Found != tt.want {
			t.Fatalf("%s failed", tt.path)
		}

		if m.Found && m.Path != tt.out {
			t.Fatalf("expected %s got %s", tt.out, m.Path)
		}
	}
}
