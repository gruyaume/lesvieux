package server_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruyaume/lesvieux/internal/db"
	"github.com/gruyaume/lesvieux/internal/server"
)

func TestNewSuccess(t *testing.T) {
	dbQueries, err := db.Initialize(":memory:")
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}

	certPath := filepath.Join("testdata", "cert.pem")
	keyPath := filepath.Join("testdata", "key.pem")
	cert, err := os.ReadFile(certPath)
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}
	_, err = server.New(1234, []byte(cert), []byte(key), dbQueries)
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}
}
