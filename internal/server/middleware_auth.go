package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gruyaume/lesvieux/internal/db"
)

const (
	UserRole  int64 = 0
	AdminRole int64 = 1
)

type contextKey string

const userIDKey = contextKey("userID")

// The adminOnly middleware checks if the user has admin role before allowing access to the handler.
func adminOnly(jwtSecret []byte, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromAuthorizationHeader(r.Header.Get("Authorization"), jwtSecret)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "auth failed: %s", err)
			return
		}

		if claims.Role != AdminRole {
			writeError(w, http.StatusForbidden, "forbidden: admin access required")
			return
		}

		// Set the user ID in the request context for further handlers
		ctxWithUserID := context.WithValue(r.Context(), userIDKey, claims.ID)
		r = r.WithContext(ctxWithUserID)

		handler(w, r)
	}
}

func Me(jwtSecret []byte, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromAuthorizationHeader(r.Header.Get("Authorization"), jwtSecret)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "auth failed: %s", err)
			return
		}

		// Set the user ID in the request context for further handlers
		ctxWithUserID := context.WithValue(r.Context(), userIDKey, claims.ID)
		r = r.WithContext(ctxWithUserID)

		handler(w, r)
	}
}

// The adminOrFirstUser middleware checks if the user has admin role or if the user is the first user before allowing access to the handler.
func adminOrFirstUser(jwtSecret []byte, db *db.Queries, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		numUsers, err := db.NumAccounts(context.Background())
		if err != nil {
			log.Println("couldn't retrieve accounts: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		if numUsers > 0 {
			claims, err := getClaimsFromAuthorizationHeader(r.Header.Get("Authorization"), jwtSecret)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "auth failed: %s", err)
				return
			}

			if claims.Role != AdminRole {
				writeError(w, http.StatusForbidden, "forbidden: admin access required")
				return
			}

			ctxWithUserID := context.WithValue(r.Context(), userIDKey, claims.ID)
			r = r.WithContext(ctxWithUserID)
		}

		handler(w, r)
	}
}

func getClaimsFromAuthorizationHeader(header string, jwtSecret []byte) (*jwtLesVieuxClaims, error) {
	if header == "" {
		return nil, fmt.Errorf("authorization header not found")
	}
	bearerToken := strings.Split(header, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return nil, fmt.Errorf("authorization header couldn't be processed. The expected format is 'Bearer <token>'")
	}
	claims, err := getClaimsFromJWT(bearerToken[1], jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("token is not valid: %s", err)
	}
	return claims, nil
}

func getClaimsFromJWT(bearerToken string, jwtSecret []byte) (*jwtLesVieuxClaims, error) {
	claims := jwtLesVieuxClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}
