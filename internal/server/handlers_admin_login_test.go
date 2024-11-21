package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
)

type AdminLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginResponseResult struct {
	Token string `json:"token"`
}

type AdminLoginResponse struct {
	Result AdminLoginResponseResult `json:"result"`
	Error  string                   `json:"error,omitempty"`
}

func adminLogin(url string, client *http.Client, data *AdminLoginParams) (int, *AdminLoginResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/admin/login", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var loginResponse AdminLoginResponse
	if err := json.NewDecoder(res.Body).Decode(&loginResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &loginResponse, nil
}

func TestAdminLoginEndToEnd(t *testing.T) {
	ts, config, err := setupServer()
	if err != nil {
		t.Fatalf("couldn't create test server: %s", err)
	}
	defer ts.Close()
	client := ts.Client()

	t.Run("Create admin user", func(t *testing.T) {
		adminUser := &CreateAdminAccountParams{
			Email:    "testadmin",
			Password: "Admin123",
		}
		statusCode, _, err := createAdminAccount(ts.URL, client, "", adminUser)
		if err != nil {
			t.Fatalf("couldn't create admin user: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}
	})

	t.Run("Login success", func(t *testing.T) {
		adminUser := &AdminLoginParams{
			Email:    "testadmin",
			Password: "Admin123",
		}
		statusCode, loginResponse, err := adminLogin(ts.URL, client, adminUser)
		if err != nil {
			t.Fatalf("couldn't login admin user: %s", err)
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
			if claims["email"] != "testadmin" {
				t.Fatalf("expected email %q, got %q", "testadmin", claims["email"])
			}
		} else {
			t.Fatalf("invalid token or claims")
		}
	})

	t.Run("Login failure missing username", func(t *testing.T) {
		invalidUser := &AdminLoginParams{
			Email:    "",
			Password: "Admin123",
		}
		statusCode, loginResponse, err := adminLogin(ts.URL, client, invalidUser)
		if err != nil {
			t.Fatalf("couldn't login admin user: %s", err)
		}
		if statusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, statusCode)
		}
		if loginResponse.Error != "Email is required" {
			t.Fatalf("expected error %q, got %q", "Email is required", loginResponse.Error)
		}
	})

	t.Run("Login failure missing password", func(t *testing.T) {
		invalidUser := &AdminLoginParams{
			Email:    "testadmin",
			Password: "",
		}
		statusCode, loginResponse, err := adminLogin(ts.URL, client, invalidUser)
		if err != nil {
			t.Fatalf("couldn't login admin user: %s", err)
		}
		if statusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, statusCode)
		}
		if loginResponse.Error != "Password is required" {
			t.Fatalf("expected error %q, got %q", "Password is required", loginResponse.Error)
		}
	})

	t.Run("Login failure invalid password", func(t *testing.T) {
		invalidUser := &AdminLoginParams{
			Email:    "testadmin",
			Password: "a-wrong-password",
		}
		statusCode, loginResponse, err := adminLogin(ts.URL, client, invalidUser)
		if err != nil {
			t.Fatalf("couldn't login admin user: %s", err)
		}
		if statusCode != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, statusCode)
		}

		if loginResponse.Error != "The username or password is incorrect. Try again." {
			t.Fatalf("expected error %q, got %q", "The username or password is incorrect. Try again.", loginResponse.Error)
		}
	})

	t.Run("Login failure invalid username", func(t *testing.T) {
		invalidUser := &AdminLoginParams{
			Email:    "not-existing-user",
			Password: "Admin123",
		}
		statusCode, loginResponse, err := adminLogin(ts.URL, client, invalidUser)
		if err != nil {
			t.Fatalf("couldn't login admin user: %s", err)
		}
		if statusCode != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, statusCode)
		}

		if loginResponse.Error != "The username or password is incorrect. Try again." {
			t.Fatalf("expected error %q, got %q", "The username or password is incorrect. Try again.", loginResponse.Error)
		}
	})
}
