package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
)

type EmployerLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmployerLoginResponseResult struct {
	Token string `json:"token"`
}

type EmployerLoginResponse struct {
	Result EmployerLoginResponseResult `json:"result"`
	Error  string                      `json:"error,omitempty"`
}

func employerLogin(url string, client *http.Client, data *EmployerLoginParams) (int, *EmployerLoginResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/employers/login", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var loginResponse EmployerLoginResponse
	if err := json.NewDecoder(res.Body).Decode(&loginResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &loginResponse, nil
}

func TestEmployerLoginEndToEnd(t *testing.T) {
	ts, config, err := setupServer()
	if err != nil {
		t.Fatalf("couldn't create test server: %s", err)
	}
	defer ts.Close()
	client := ts.Client()

	var token string
	t.Run("prepare user accounts and tokens", prepareAdminAccount(ts.URL, client, &token))

	t.Run("Create employer", func(t *testing.T) {
		employer := &CreateEmployerParams{
			Name: "Test Employer",
		}
		statusCode, _, err := createEmployer(ts.URL, client, token, employer)
		if err != nil {
			t.Fatalf("couldn't create employer: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}
	})

	t.Run("Create employer user", func(t *testing.T) {
		employerUser := &CreateEmployerAccountParams{
			Email:    "testemployer",
			Password: "Employer123!",
		}
		employerId := "1"
		statusCode, _, err := createEmployerAccount(ts.URL, client, token, employerId, employerUser)
		if err != nil {
			t.Fatalf("couldn't create employer user: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}
	})

	t.Run("Login success", func(t *testing.T) {
		employerUser := &EmployerLoginParams{
			Email:    "testemployer",
			Password: "Employer123!",
		}
		statusCode, loginResponse, err := employerLogin(ts.URL, client, employerUser)
		if err != nil {
			t.Fatalf("couldn't login employer user: %s", err)
		}
		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}
		if loginResponse.Result.Token == "" {
			t.Fatalf("expected token, got empty string")
		}
		token, err := jwt.Parse(loginResponse.Result.Token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			t.Fatalf("couldn't parse token: %s", err)
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["email"] != "testemployer" {
				t.Fatalf("expected email %q, got %q", "testemployer", claims["email"])
			}
		} else {
			t.Fatalf("invalid token or claims")
		}
	})

	t.Run("Login failure missing username", func(t *testing.T) {
		invalidUser := &EmployerLoginParams{
			Email:    "",
			Password: "Employer123",
		}
		statusCode, loginResponse, err := employerLogin(ts.URL, client, invalidUser)
		if err != nil {
			t.Fatalf("couldn't login employer user: %s", err)
		}
		if statusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, statusCode)
		}
		if loginResponse.Error != "Email is required" {
			t.Fatalf("expected error %q, got %q", "Email is required", loginResponse.Error)
		}
	})

	t.Run("Login failure missing password", func(t *testing.T) {
		invalidUser := &EmployerLoginParams{
			Email:    "testemployer",
			Password: "",
		}
		statusCode, loginResponse, err := employerLogin(ts.URL, client, invalidUser)
		if err != nil {
			t.Fatalf("couldn't login employer user: %s", err)
		}
		if statusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, statusCode)
		}
		if loginResponse.Error != "Password is required" {
			t.Fatalf("expected error %q, got %q", "Password is required", loginResponse.Error)
		}
	})

	t.Run("Login failure invalid password", func(t *testing.T) {
		invalidUser := &EmployerLoginParams{
			Email:    "testemployer",
			Password: "a-wrong-password",
		}
		statusCode, loginResponse, err := employerLogin(ts.URL, client, invalidUser)
		if err != nil {
			t.Fatalf("couldn't login employer user: %s", err)
		}
		if statusCode != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, statusCode)
		}

		if loginResponse.Error != "The email or password is incorrect. Try again." {
			t.Fatalf("expected error %q, got %q", "The email or password is incorrect. Try again.", loginResponse.Error)
		}
	})

	t.Run("Login failure invalid email", func(t *testing.T) {
		invalidUser := &EmployerLoginParams{
			Email:    "not-existing-user",
			Password: "Employer123",
		}
		statusCode, loginResponse, err := employerLogin(ts.URL, client, invalidUser)
		if err != nil {
			t.Fatalf("couldn't login employer user: %s", err)
		}
		if statusCode != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, statusCode)
		}

		if loginResponse.Error != "The email or password is incorrect. Try again." {
			t.Fatalf("expected error %q, got %q", "The email or password is incorrect. Try again.", loginResponse.Error)
		}
	})
}
