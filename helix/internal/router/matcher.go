package router

import (
	"strings"

	"github.com/shubhkr72/helix/internal/config"
)

type Match struct {
	Route config.Route
	Path  string
	Found bool
}

func MatchRoute(routes []config.Route, path string) Match {

	var best config.Route
	bestLen := -1

	for _, r := range routes {

		if path == r.Path || strings.HasPrefix(path, r.Path+"/") {

			if len(r.Path) > bestLen {
				best = r
				bestLen = len(r.Path)
			}
		}
	}

	if bestLen == -1 {
		return Match{}
	}

	newPath := path

	if best.StripPrefix {

		newPath = strings.TrimPrefix(path, best.Path)

		if newPath == "" {
			newPath = "/"
		}
	}

	return Match{
		Route: best,
		Path:  newPath,
		Found: true,
	}
}
