package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gruyaume/lesvieux/internal/db"
)

type CreateAdminAccountParams struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAdminAccountResponse struct {
	ID int64 `json:"id"`
}

type GetAdminAccountResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type ChangeAdminAccountPasswordParams struct {
	Password string `json:"password"`
}

type ChangeAdminAccountPasswordResponse struct {
	ID int64 `json:"id"`
}

func ListAdminAccounts(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := env.DBQueries.ListAdminAccounts(context.Background())
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		accountsResponse := make([]GetAdminAccountResponse, 0, len(accounts))
		for i := range accounts {
			accountsResponse = append(accountsResponse, GetAdminAccountResponse{
				ID:    accounts[i].ID,
				Email: accounts[i].Email,
			})
		}
		err = writeJSON(w, accountsResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// GetAdminAccount receives an id as a path parameter, and
// returns the corresponding AdminAccount
func GetAdminAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("account_id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		var DBAdminAccount db.AdminAccount
		DBAdminAccount, err = env.DBQueries.GetAdminAccount(context.Background(), idInt)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Admin Account not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		accountResponse := GetAdminAccountResponse{
			ID:    DBAdminAccount.ID,
			Email: DBAdminAccount.Email,
		}
		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, accountResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func GetMyAdminAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromAuthorizationHeader(r.Header.Get("Authorization"), env.JWTSecret)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		DBAdminAccount, err := env.DBQueries.GetAdminAccountByEmail(context.Background(), claims.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Admin Account not found")
				return
			}
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		accountResponse := GetAdminAccountResponse{
			ID:    DBAdminAccount.ID,
			Email: DBAdminAccount.Email,
		}
		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, accountResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// CreateAdminAccount creates a new AdminAccount, and returns the id of the created row
func CreateAdminAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var account CreateAdminAccountParams
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if account.Email == "" {
			writeError(w, http.StatusBadRequest, "Email is required")
			return
		}
		if account.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		if !validatePassword(account.Password) {
			writeError(
				w,
				http.StatusBadRequest,
				"Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol.",
			)
			return
		}

		passwordHash, err := GeneratePasswordHash(account.Password)
		if err != nil {
			log.Println("Failed to generate password hash: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		newAdminAccountParams := db.CreateAdminAccountParams{
			Email:        account.Email,
			PasswordHash: passwordHash,
		}
		newAdminAccount, err := env.DBQueries.CreateAdminAccount(context.Background(), newAdminAccountParams)
		if err != nil {
			log.Println("Failed to create account: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusCreated)
		response := CreateAdminAccountResponse{ID: newAdminAccount.ID}
		err = writeJSON(w, response)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// DeleteAdminAccount handler receives an id as a path parameter,
// deletes the corresponding AdminAccount, and returns a http.StatusNoContent on success
func DeleteAdminAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("account_id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		_, err = env.DBQueries.GetAdminAccount(context.Background(), idInt)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Admin Account not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		err = env.DBQueries.DeleteAdminAccount(context.Background(), idInt)
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusAccepted)
		response := map[string]any{"id": idInt}
		err = writeJSON(w, response)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func ChangeAdminAccountPassword(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("account_id")
		var changeAdminAccountPassword ChangeAdminAccountPasswordParams
		if err := json.NewDecoder(r.Body).Decode(&changeAdminAccountPassword); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		pathIDInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}

		_, err = env.DBQueries.GetAdminAccount(context.Background(), pathIDInt64)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Admin Account not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if changeAdminAccountPassword.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		if !validatePassword(changeAdminAccountPassword.Password) {
			writeError(
				w,
				http.StatusBadRequest,
				"Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol.",
			)
			return
		}
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		passwordHash, err := GeneratePasswordHash(changeAdminAccountPassword.Password)
		if err != nil {
			log.Println("Failed to generate password hash: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		updateAdminAccountParams := db.UpdateAdminAccountParams{
			ID:           idInt,
			PasswordHash: passwordHash,
		}
		err = env.DBQueries.UpdateAdminAccount(context.Background(), updateAdminAccountParams)
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusOK)
		changeAdminAccountPasswordResponse := ChangeAdminAccountPasswordResponse{
			ID: idInt,
		}
		err = writeJSON(w, changeAdminAccountPasswordResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func ChangeMyAdminAccountPassword(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromAuthorizationHeader(r.Header.Get("Authorization"), env.JWTSecret)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		account, err := env.DBQueries.GetAdminAccountByEmail(context.Background(), claims.Email)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		id := strconv.FormatInt(account.ID, 10)
		var changeAdminAccountPassword ChangeAdminAccountPasswordParams
		if err := json.NewDecoder(r.Body).Decode(&changeAdminAccountPassword); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if changeAdminAccountPassword.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		if !validatePassword(changeAdminAccountPassword.Password) {
			writeError(
				w,
				http.StatusBadRequest,
				"Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol.",
			)
			return
		}
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		passwordHash, err := GeneratePasswordHash(changeAdminAccountPassword.Password)
		if err != nil {
			log.Println("Failed to generate password hash: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		updateAdminAccountParams := db.UpdateAdminAccountParams{
			ID:           idInt,
			PasswordHash: passwordHash,
		}
		err = env.DBQueries.UpdateAdminAccount(context.Background(), updateAdminAccountParams)
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusOK)
		changeAdminAccountPasswordResponse := ChangeAdminAccountPasswordResponse{
			ID: idInt,
		}
		err = writeJSON(w, changeAdminAccountPasswordResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}
