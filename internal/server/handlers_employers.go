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

type CreateEmployerParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CreateEmployerResponse struct {
	ID int64 `json:"id"`
}

type GetEmployerResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ListEmployers(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employers, err := env.DBQueries.ListEmployers(context.Background())
		if err != nil {
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		employersResponse := make([]GetEmployerResponse, 0, len(employers))
		for i := range employers {
			employersResponse = append(employersResponse, GetEmployerResponse{
				ID:   employers[i].ID,
				Name: employers[i].Name,
			})
		}
		err = writeJSON(w, employersResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// GetEmployer receives an id as a path parameter, and
// returns the corresponding Employer
func GetEmployer(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("employer_id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		var DBEmployer db.Employer
		DBEmployer, err = env.DBQueries.GetEmployer(context.Background(), idInt)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Employer not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		employerResponse := GetEmployerResponse{
			ID:   DBEmployer.ID,
			Name: DBEmployer.Name,
		}
		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, employerResponse)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// CreateEmployer creates a new Employer, and returns the id of the created row
func CreateEmployer(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employer CreateEmployerParams
		if err := json.NewDecoder(r.Body).Decode(&employer); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if employer.Name == "" {
			writeError(w, http.StatusBadRequest, "Name is required")
			return
		}

		newEmployer, err := env.DBQueries.CreateEmployer(context.Background(), employer.Name)
		if err != nil {
			log.Println("Failed to create employer: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		w.WriteHeader(http.StatusCreated)
		response := CreateEmployerResponse{ID: newEmployer.ID}
		err = writeJSON(w, response)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}

// DeleteEmployer handler receives an id as a path parameter,
// deletes the corresponding Employer, and returns a http.StatusNoContent on success
func DeleteEmployer(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("employer_id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid id")
			return
		}
		_, err = env.DBQueries.GetEmployer(context.Background(), idInt)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "Employer not found")
				return
			}
			log.Println(err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		err = env.DBQueries.DeleteEmployer(context.Background(), idInt)
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
