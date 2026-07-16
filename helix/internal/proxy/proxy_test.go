package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newBackend() *httptest.Server {

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resp := map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
			"query":  r.URL.RawQuery,
			"xff":    r.Header.Get("X-Forwarded-For"),
			"body":   "",
		}

		if r.Body != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			resp["body"] = buf.String()
		}

		json.NewEncoder(w).Encode(resp)
	}))
}

func TestGET(t *testing.T) {

	backend := newBackend()
	defer backend.Close()

	p, _ := New(backend.URL)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users?id=1",
		nil,
	)

	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatal()
	}
}

func TestPOSTBody(t *testing.T) {

	backend := newBackend()
	defer backend.Close()

	p, _ := New(backend.URL)

	req := httptest.NewRequest(
		http.MethodPost,
		"/users",
		bytes.NewBufferString("hello"),
	)

	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatal()
	}
}

func TestQuery(t *testing.T) {

	backend := newBackend()
	defer backend.Close()

	p, _ := New(backend.URL)

	req := httptest.NewRequest(
		http.MethodGet,
		"/test?a=1&b=2",
		nil,
	)

	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatal()
	}
}

func TestHeaderForwarding(t *testing.T) {

	backend := newBackend()
	defer backend.Close()

	p, _ := New(backend.URL)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set("User-Agent", "Go-Test")

	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatal()
	}
}

func TestSpoofedXFF(t *testing.T) {

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("X-Forwarded-For") == "1.2.3.4" {
			t.Fatal("spoofed header forwarded")
		}

	}))
	defer backend.Close()

	p, _ := New(backend.URL)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set("X-Forwarded-For", "1.2.3.4")

	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)
}

func TestBackendDown(t *testing.T) {

	p, _ := New("http://127.0.0.1:9999")

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatal()
	}
}

func TestContextCancel(t *testing.T) {

	backend := newBackend()
	defer backend.Close()

	p, _ := New(backend.URL)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	ctx, cancel := context.WithCancel(req.Context())
	cancel()

	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req.WithContext(ctx))
}
