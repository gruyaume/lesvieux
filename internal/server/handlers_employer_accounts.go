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

type CreateEmployerAccountParams struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Role     int64  `json:"role"`
	Password string `json:"password"`
}

type CreateEmployerAccountResponse struct {
	ID int64 `json:"id"`
}

type GetEmployerAccountResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  int64  `json:"role"`
}

type ChangeEmployerAccountPasswordParams struct {
	Password string `json:"password"`
}

type ChangeEmployerAccountPasswordResponse struct {
	ID int64 `json:"id"`
}

func ListEmployerAccounts(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := env.DBQueries.ListEmployerAccounts(context.Background())
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		accountsResponse := make([]GetEmployerAccountResponse, 0, len(accounts))
		for i := range accounts {
			accountsResponse = append(accountsResponse, GetEmployerAccountResponse{
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

// GetEmployerAccount receives an id as a path parameter, and
// returns the corresponding EmployerAccount
func GetEmployerAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("user_id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		var DBEmployerAccount db.EmployerAccount
		DBEmployerAccount, err = env.DBQueries.GetEmployerAccount(context.Background(), idInt)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "EmployerAccount not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		accountResponse := GetEmployerAccountResponse{
			ID:    DBEmployerAccount.ID,
			Email: DBEmployerAccount.Email,
		}
		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, accountResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func GetMyEmployerAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromAuthorizationHeader(r.Header.Get("Authorization"), env.JWTSecret)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		DBEmployerAccount, err := env.DBQueries.GetEmployerAccountByEmail(context.Background(), claims.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "EmployerAccount not found")
				return
			}
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		accountResponse := GetEmployerAccountResponse{
			ID:    DBEmployerAccount.ID,
			Email: DBEmployerAccount.Email,
		}
		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, accountResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// CreateEmployerAccount creates a new EmployerAccount, and returns the id of the created row
func CreateEmployerAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var account CreateEmployerAccountParams
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

		_, err := env.DBQueries.GetEmployerAccountByEmail(context.Background(), account.Email)
		if err == nil {
			writeError(w, http.StatusConflict, "Account already exists")
			return
		}

		passwordHash, err := GeneratePasswordHash(account.Password)
		if err != nil {
			log.Println("Failed to generate password hash: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		newEmployerAccountParams := db.CreateEmployerAccountParams{
			Email:        account.Email,
			PasswordHash: passwordHash,
		}
		newEmployerAccount, err := env.DBQueries.CreateEmployerAccount(context.Background(), newEmployerAccountParams)
		if err != nil {
			log.Println("Failed to create account: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusCreated)
		response := CreateEmployerAccountResponse{ID: newEmployerAccount.ID}
		err = writeJSON(w, response)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// DeleteEmployerAccount handler receives an id as a path parameter,
// deletes the corresponding EmployerAccount, and returns a http.StatusNoContent on success
func DeleteEmployerAccount(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("user_id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		_, err = env.DBQueries.GetEmployerAccount(context.Background(), idInt)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "EmployerAccount not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		err = env.DBQueries.DeleteEmployerAccount(context.Background(), idInt)
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

func ChangeEmployerAccountPassword(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("user_id")
		var changeEmployerAccountPassword ChangeEmployerAccountPasswordParams
		if err := json.NewDecoder(r.Body).Decode(&changeEmployerAccountPassword); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		pathIDInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}

		_, err = env.DBQueries.GetEmployerAccount(context.Background(), pathIDInt64)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "EmployerAccount not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if changeEmployerAccountPassword.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		if !validatePassword(changeEmployerAccountPassword.Password) {
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
		passwordHash, err := GeneratePasswordHash(changeEmployerAccountPassword.Password)
		if err != nil {
			log.Println("Failed to generate password hash: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		updateEmployerAccountParams := db.UpdateEmployerAccountParams{
			ID:           idInt,
			PasswordHash: passwordHash,
		}
		err = env.DBQueries.UpdateEmployerAccount(context.Background(), updateEmployerAccountParams)
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusOK)
		changeEmployerAccountPasswordResponse := ChangeEmployerAccountPasswordResponse{
			ID: idInt,
		}
		err = writeJSON(w, changeEmployerAccountPasswordResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

func ChangeMyEmployerAccountPassword(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromAuthorizationHeader(r.Header.Get("Authorization"), env.JWTSecret)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		account, err := env.DBQueries.GetEmployerAccountByEmail(context.Background(), claims.Email)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		id := strconv.FormatInt(account.ID, 10)
		var changeEmployerAccountPassword ChangeEmployerAccountPasswordParams
		if err := json.NewDecoder(r.Body).Decode(&changeEmployerAccountPassword); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if changeEmployerAccountPassword.Password == "" {
			writeError(w, http.StatusBadRequest, "Password is required")
			return
		}
		if !validatePassword(changeEmployerAccountPassword.Password) {
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
		passwordHash, err := GeneratePasswordHash(changeEmployerAccountPassword.Password)
		if err != nil {
			log.Println("Failed to generate password hash: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		updateEmployerAccountParams := db.UpdateEmployerAccountParams{
			ID:           idInt,
			PasswordHash: passwordHash,
		}
		err = env.DBQueries.UpdateEmployerAccount(context.Background(), updateEmployerAccountParams)
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusOK)
		changeEmployerAccountPasswordResponse := ChangeEmployerAccountPasswordResponse{
			ID: idInt,
		}
		err = writeJSON(w, changeEmployerAccountPasswordResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}
