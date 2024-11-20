package server

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gruyaume/lesvieux/internal/db"
)

type HandlerConfig struct {
	DBQueries *db.Queries
	JWTSecret []byte
}

func generateJWTSecret() ([]byte, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return bytes, fmt.Errorf("failed to generate JWT secret: %w", err)
	}
	return bytes, nil
}

func New(port int, cert []byte, key []byte, dbQueries *db.Queries) (*http.Server, error) {
	jwtSecret, err := generateJWTSecret()
	if err != nil {
		return nil, err
	}
	env := &HandlerConfig{
		DBQueries: dbQueries,
		JWTSecret: jwtSecret,
	}
	router := NewHebdoRouter(env)

	serverCerts, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Addr: fmt.Sprintf(":%d", port),

		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{serverCerts},
		},
	}

	return s, nil
}
