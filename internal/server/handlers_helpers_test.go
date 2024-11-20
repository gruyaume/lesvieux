// Contains helper functions for testing the server
package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruyaume/lesvieux/internal/db"
	"github.com/gruyaume/lesvieux/internal/server"
)

var adminUser = CreateAccountParams{
	Username: "testadmin",
	Password: "Admin123",
}

var validUser = CreateAccountParams{
	Username: "testuser",
	Password: "userPass!",
}

func setupServer() (*httptest.Server, *server.HandlerConfig, error) {
	dbQueries, err := db.Initialize(":memory:")
	if err != nil {
		return nil, nil, err
	}
	config := &server.HandlerConfig{
		DBQueries: dbQueries,
	}
	ts := httptest.NewTLSServer(server.NewLesVieuxRouter(config))
	return ts, config, nil
}

func prepareUserAccounts(url string, client *http.Client, adminToken, nonAdminToken *string) func(*testing.T) {
	return func(t *testing.T) {
		statusCode, _, err := createAccount(url, client, "", &adminUser)
		if err != nil {
			t.Fatalf("couldn't create admin user: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("creating the first request should succeed when unauthorized. status code received: %d", statusCode)
		}

		loginParams := LoginParams{
			Username: adminUser.Username,
			Password: adminUser.Password,
		}
		statusCode, loginResponse, err := login(url, client, &loginParams)
		if err != nil {
			t.Fatalf("couldn't login admin user: %s", err)
		}
		if statusCode != http.StatusOK {
			t.Fatalf("the admin login request should have succeeded. status code received: %d", statusCode)
		}

		*adminToken = loginResponse.Result.Token

		statusCode, _, err = createAccount(url, client, *adminToken, &validUser)
		if err != nil {
			t.Fatalf("couldn't create test user: %s", err)
		}

		if statusCode != http.StatusCreated {
			t.Fatalf("creating the second request should have succeeded given the admin auth header. status code received: %d", statusCode)
		}

		loginParams = LoginParams{
			Username: validUser.Username,
			Password: validUser.Password,
		}
		statusCode, loginResponse, err = login(url, client, &loginParams)
		if err != nil {
			t.Fatalf("couldn't login test user: %s", err)
		}
		if statusCode != http.StatusOK {
			t.Fatalf("the test user login request should have succeeded. status code received: %d", statusCode)
		}

		*nonAdminToken = loginResponse.Result.Token
	}
}
