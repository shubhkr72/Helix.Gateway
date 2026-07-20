package proxy

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/shubhkr72/helix/internal/auth"
)

type Gateway struct {
	proxy *httputil.ReverseProxy
}

func New(target string) (*Gateway, error) {

	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	rp := &httputil.ReverseProxy{
		Transport: NewTransport(),

		Rewrite: func(pr *httputil.ProxyRequest) {

			pr.Out.Header.Del("X-Forwarded-For")
			pr.Out.Header.Del("X-Forwarded-Host")
			pr.Out.Header.Del("X-Forwarded-Proto")

			pr.Out.Header.Del("X-User-ID")
			pr.Out.Header.Del("X-Email")
			pr.Out.Header.Del("X-Roles")

			pr.SetURL(u)
			pr.SetXForwarded()

			pr.Out.Host = u.Host

			principal := auth.GetPrincipal(pr.In.Context())
			if principal == nil {
				return
			}

			pr.Out.Header.Set("X-User-ID", principal.UserID)
			pr.Out.Header.Set("X-Email", principal.Email)
			pr.Out.Header.Set("X-Roles", strings.Join(principal.Roles, ","))
		},

		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {

			id := requestID()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)

			fmt.Fprintf(
				w,
				`{"error":"bad gateway","request_id":"%s"}`,
				id,
			)
		},
	}

	return &Gateway{
		proxy: rp,
	}, nil
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	g.proxy.ServeHTTP(
		w,
		r.WithContext(ctx),
	)
}

func requestID() string {

	b := make([]byte, 8)

	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}
