package config

import "testing"

func TestDuplicateID(t *testing.T) {

	cfg := Config{
		Routes: []Route{
			{ID: "a", Path: "/a", Backend: []string{"http://localhost"}},
			{ID: "a", Path: "/b", Backend: []string{"http://localhost"}},
		},
	}

	if validate(cfg) == nil {
		t.Fatal("expected error")
	}
}

func TestEmptyBackend(t *testing.T) {

	cfg := Config{
		Routes: []Route{
			{ID: "a", Path: "/a"},
		},
	}

	if validate(cfg) == nil {
		t.Fatal("expected error")
	}
}
