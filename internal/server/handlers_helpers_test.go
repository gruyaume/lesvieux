// Contains helper functions for testing the server
package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruyaume/lesvieux/internal/db"
	"github.com/gruyaume/lesvieux/internal/server"
)

var adminUser = CreateAdminAccountParams{
	Email:    "testadmin",
	Password: "Admin123",
}

var validEmployerAccount = CreateEmployerAccountParams{
	Email:    "testemployer",
	Password: "Employerpass123!",
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

func prepareAdminAccount(url string, client *http.Client, token *string) func(*testing.T) {
	return func(t *testing.T) {
		statusCode, _, err := createAdminAccount(url, client, "", &adminUser)
		if err != nil {
			t.Fatalf("couldn't create employer account: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("creating the first request should succeed when unauthorized. status code received: %d", statusCode)
		}

		loginParams := AdminLoginParams{
			Email:    adminUser.Email,
			Password: adminUser.Password,
		}
		statusCode, loginResponse, err := adminLogin(url, client, &loginParams)
		if err != nil {
			t.Fatalf("couldn't login admin user: %s", err)
		}
		if statusCode != http.StatusOK {
			t.Fatalf("the admin login request should have succeeded. status code received: %d", statusCode)
		}

		*token = loginResponse.Result.Token

	}
}

func prepareEmployerAccount(url string, client *http.Client, adminToken *string, employerToken *string) func(*testing.T) {
	return func(t *testing.T) {
		statusCode, _, err := createEmployerAccount(url, client, *adminToken, &validEmployerAccount)
		if err != nil {
			t.Fatalf("couldn't create employer account: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("creating the first request should succeed when unauthorized. status code received: %d", statusCode)
		}

		loginParams := EmployerLoginParams{
			Email:    validEmployerAccount.Email,
			Password: validEmployerAccount.Password,
		}
		statusCode, loginResponse, err := employerLogin(url, client, &loginParams)
		if err != nil {
			t.Fatalf("couldn't login employer user: %s", err)
		}
		if statusCode != http.StatusOK {
			t.Fatalf("the employer login request should have succeeded. status code received: %d", statusCode)
		}

		*employerToken = loginResponse.Result.Token

	}
}
