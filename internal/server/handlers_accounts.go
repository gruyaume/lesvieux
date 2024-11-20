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

type CreateAccountParams struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     int64  `json:"role"`
	Password string `json:"password"`
}

type CreateAccountResponse struct {
	ID int64 `json:"id"`
}

type GetAccountResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     int64  `json:"role"`
}

type ChangeAccountPasswordParams struct {
	Password string `json:"password"`
}

type ChangeAccountPasswordResponse struct {
	ID int64 `json:"id"`
}

func ListAccounts(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := env.DBQueries.ListAccounts(context.Background())
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		accountsResponse := make([]GetAccountResponse, 0, len(accounts))
		for i := range accounts {
			accountsResponse = append(accountsResponse, GetAccountResponse{
				ID:       accounts[i].ID,
				Username: accounts[i].Username,
				Role:     accounts[i].Role,
			})
		}
		err = writeJSON(w, accountsResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// GetAccount receives an id as a path parameter, and
// returns the corresponding Account
func GetAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("user_id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		var DBAccount db.Account
		DBAccount, err = env.DBQueries.GetAccount(context.Background(), idInt)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Account not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		accountResponse := GetAccountResponse{
			ID:       DBAccount.ID,
			Username: DBAccount.Username,
			Role:     DBAccount.Role,
		}
		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, accountResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func GetMyAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromAuthorizationHeader(r.Header.Get("Authorization"), env.JWTSecret)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		DBAccount, err := env.DBQueries.GetAccountByUsername(context.Background(), claims.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Account not found")
				return
			}
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		accountResponse := GetAccountResponse{
			ID:       DBAccount.ID,
			Username: DBAccount.Username,
			Role:     DBAccount.Role,
		}
		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, accountResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// CreateAccount creates a new Account, and returns the id of the created row
func CreateAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var account CreateAccountParams
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if account.Username == "" {
			writeError(w, http.StatusBadRequest, "Username is required")
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
		numAccounts, err := env.DBQueries.NumAccounts(context.Background())
		if err != nil {
			log.Println("Failed to retrieve accounts: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		role := UserRole
		if numAccounts == 0 {
			role = AdminRole
		}
		passwordHash, err := GeneratePasswordHash(account.Password)
		if err != nil {
			log.Println("Failed to generate password hash: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		newAccountParams := db.CreateAccountParams{
			Username:     account.Username,
			PasswordHash: passwordHash,
			Role:         role,
		}
		newAccount, err := env.DBQueries.CreateAccount(context.Background(), newAccountParams)
		if err != nil {
			log.Println("Failed to create account: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusCreated)
		response := CreateAccountResponse{ID: newAccount.ID}
		err = writeJSON(w, response)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// DeleteAccount handler receives an id as a path parameter,
// deletes the corresponding Account, and returns a http.StatusNoContent on success
func DeleteAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("user_id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		account, err := env.DBQueries.GetAccount(context.Background(), idInt)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Account not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if account.Role == AdminRole {
			writeError(w, http.StatusBadRequest, "deleting an Admin account is not allowed.")
			return
		}
		err = env.DBQueries.DeleteAccount(context.Background(), idInt)
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

func ChangeAccountPassword(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("user_id")
		var changeAccountPassword ChangeAccountPasswordParams
		if err := json.NewDecoder(r.Body).Decode(&changeAccountPassword); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		pathIDInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}

		_, err = env.DBQueries.GetAccount(context.Background(), pathIDInt64)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Account not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if changeAccountPassword.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		if !validatePassword(changeAccountPassword.Password) {
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
		passwordHash, err := GeneratePasswordHash(changeAccountPassword.Password)
		if err != nil {
			log.Println("Failed to generate password hash: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		updateAccountParams := db.UpdateAccountParams{
			ID:           idInt,
			PasswordHash: passwordHash,
		}
		err = env.DBQueries.UpdateAccount(context.Background(), updateAccountParams)
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusOK)
		changeAccountPasswordResponse := ChangeAccountPasswordResponse{
			ID: idInt,
		}
		err = writeJSON(w, changeAccountPasswordResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func ChangeMyAccountPassword(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromAuthorizationHeader(r.Header.Get("Authorization"), env.JWTSecret)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		account, err := env.DBQueries.GetAccountByUsername(context.Background(), claims.Username)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		id := strconv.FormatInt(account.ID, 10)
		var changeAccountPassword ChangeAccountPasswordParams
		if err := json.NewDecoder(r.Body).Decode(&changeAccountPassword); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if changeAccountPassword.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		if !validatePassword(changeAccountPassword.Password) {
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
		passwordHash, err := GeneratePasswordHash(changeAccountPassword.Password)
		if err != nil {
			log.Println("Failed to generate password hash: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		updateAccountParams := db.UpdateAccountParams{
			ID:           idInt,
			PasswordHash: passwordHash,
		}
		err = env.DBQueries.UpdateAccount(context.Background(), updateAccountParams)
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusOK)
		changeAccountPasswordResponse := ChangeAccountPasswordResponse{
			ID: idInt,
		}
		err = writeJSON(w, changeAccountPasswordResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}
