package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type jwtAdminClaims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  int64  `json:"role"`
	jwt.StandardClaims
}

// Helper function to generate a JWT
func generateAdminJWT(id int64, email string, jwtSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtAdminClaims{
		ID:    id,
		Email: email,
		Role:  AdminRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ClaimValidity).Unix(),
		},
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AdminLogin(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest LoginParams
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if loginRequest.Email == "" {
			writeError(w, http.StatusBadRequest, "Email is required")
			return
		}
		if loginRequest.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		account, err := env.DBQueries.GetAdminAccountByEmail(context.Background(), loginRequest.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusUnauthorized, "The username or password is incorrect. Try again.")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(loginRequest.Password)); err != nil {
			writeError(w, http.StatusUnauthorized, "The username or password is incorrect. Try again.")
			return
		}
		jwt, err := generateAdminJWT(account.ID, account.Email, env.JWTSecret)
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		loginResponse := LoginResponse{
			Token: jwt,
		}
		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, loginResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}
