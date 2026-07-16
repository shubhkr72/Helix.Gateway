package proxy

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
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

			// reject spoofed forwarding headers

			pr.Out.Header.Del("X-Forwarded-For")
			pr.Out.Header.Del("X-Forwarded-Host")
			pr.Out.Header.Del("X-Forwarded-Proto")

			pr.SetURL(u)

			pr.SetXForwarded()

			pr.Out.Host = u.Host
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
	rand.Read(b)

	return hex.EncodeToString(b)
}
