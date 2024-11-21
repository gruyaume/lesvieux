package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gruyaume/lesvieux/version"
)

type GetStatusResponse struct {
	Initialized bool   `json:"initialized"`
	Version     string `json:"version"`
}

func GetStatus(env *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		numEmployerAccounts, err := env.DBQueries.NumAdminAccounts(context.Background())
		if err != nil {
			log.Println("couldn't retrieve admin accounts: " + err.Error())
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		response := GetStatusResponse{
			Initialized: numEmployerAccounts > 0,
			Version:     version.GetVersion(),
		}
		w.WriteHeader(http.StatusOK)
		err = writeJSON(w, response)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
}
