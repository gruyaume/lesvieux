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

const ClaimValidity = 1 * time.Hour

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type jwtLesVieuxClaims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     int64  `json:"role"`
	jwt.StandardClaims
}

// Helper function to generate a JWT
func generateJWT(id int64, username string, jwtSecret []byte, role int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtLesVieuxClaims{
		ID:       id,
		Username: username,
		Role:     role,
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

func Login(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest LoginParams
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if loginRequest.Username == "" {
			writeError(w, http.StatusBadRequest, "Username is required")
			return
		}
		if loginRequest.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		account, err := env.DBQueries.GetAccountByUsername(context.Background(), loginRequest.Username)
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
		jwt, err := generateJWT(account.ID, account.Username, env.JWTSecret, account.Role)
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
