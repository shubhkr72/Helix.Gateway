package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/shubhkr72/helix/internal/auth"
	"github.com/shubhkr72/helix/internal/jwt"
)

func authenticate(r *http.Request, manager *jwt.Manager) (*http.Request, error) {
	header := r.Header.Get("Authorization")

	if header == "" {
		return nil, http.ErrNoCookie
	}

	parts := strings.SplitN(header, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("invalid authorization header")
	}

	claims, err := manager.VerifyToken(parts[1])
	if err != nil {
		return nil, err
	}

	principal := &auth.Principal{
		UserID: claims.Subject,
		Email:  claims.Email,
		Roles:  claims.Roles,
	}

	ctx := auth.SetPrincipal(r.Context(), principal)

	return r.WithContext(ctx), nil
}
