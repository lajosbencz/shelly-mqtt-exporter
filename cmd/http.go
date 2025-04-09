package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newServer(cfg *config) *http.Server {
	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(collector.Registry(), promhttp.HandlerOpts{})
	server := &http.Server{
		Addr:    ":2112",
		Handler: mux,
	}
	if cfg.PrometheusUser != "" && cfg.PrometheusPass != "" {
		logger.Info("using Prometheus authentication")
		basicAuth := newBasicAuth("Prometheus - Shelly", cfg.PrometheusUser, cfg.PrometheusPass)
		mux.Handle("/metrics", basicAuth.Wrap(promHandler))
	} else {
		mux.Handle("/metrics", promHandler)
	}
	return server
}

func newBasicAuth(realm, username, password string) *basicAuth {
	return &basicAuth{
		realm:        realm,
		username:     username,
		passwordHash: sha256.Sum256([]byte(password)),
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
