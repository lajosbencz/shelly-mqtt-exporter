package main

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newServer(cfg *config) *server {
	address := fmt.Sprintf("%s:%d", cfg.PrometheusHost, cfg.PrometheusPort)
	mux := http.NewServeMux()
	httpLogger := newServerErrorLog()
	srv := &http.Server{
		Addr:         address,
		Handler:      mux,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		ErrorLog:     httpLogger,
	}

	promHandler := promhttp.HandlerFor(collector.Registry(), promhttp.HandlerOpts{})
	if cfg.PrometheusUser != "" && cfg.PrometheusPassSha != nil {
		logger.Info("using Prometheus authentication")
		basicAuth := newBasicAuth("Prometheus - Shelly", cfg.PrometheusUser, *cfg.PrometheusPassSha)
		mux.Handle("/metrics", basicAuth.Wrap(promHandler))
	} else {
		mux.Handle("/metrics", promHandler)
	}
	if cfg.IsTls() {
		logger.Info("metrics server configured with TLS")
	}
	return &server{
		Server: srv,
		config: cfg,
		mux:    mux,
	}
}

type server struct {
	*http.Server
	config *config
	mux    *http.ServeMux
}

func (s *server) Scheme() string {
	if s.config.IsTls() {
		return "https"
	}
	return "http"
}

func (s *server) ListenAndServe() error {
	if s.config.IsTls() {
		return s.Server.ListenAndServeTLS(s.config.TLSCertPath, s.config.TLSKeyPath)
	}
	return s.Server.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *server) Mux() *http.ServeMux {
	return s.mux
}

func newBasicAuth(realm, username string, passwordHash [32]byte) *basicAuth {
	return &basicAuth{
		realm:        realm,
		username:     username,
		passwordHash: passwordHash,
	}
}

type basicAuth struct {
	realm        string
	username     string
	passwordHash [32]byte
}

func (b *basicAuth) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		expectedUsernameHash := sha256.Sum256([]byte(b.username))
		expectedPasswordHash := b.passwordHash
		if !ok ||
			subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) != 1 ||
			subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) != 1 {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s", charset="UTF-8"`, b.realm))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
