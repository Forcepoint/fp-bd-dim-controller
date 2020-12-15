package auth

import (
	"encoding/json"
	"fmt"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/auth"
	"fp-dynamic-elements-manager-controller/internal/auth/structs"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"github.com/spf13/viper"
	"net/http"
)

func Login(repo *persistence.UserRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &structs.User{}
		err := json.NewDecoder(r.Body).Decode(user)

		if err != nil {
			json.NewEncoder(w).Encode(util.HttpResponse{Status: http.StatusBadRequest, Message: "Invalid request"})
			return
		}

		exception, resp := auth.FindUserAndCreateToken(user.Email, user.Password, repo)

		enc := json.NewEncoder(w)
		if exception != nil {
			w.WriteHeader(http.StatusUnauthorized)
			enc.Encode(exception)
		} else {
			w.WriteHeader(http.StatusOK)
			enc.Encode(resp)
		}
		return
	})
}

func GetRegistrationKey() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == "GET" {
			json.NewEncoder(w).Encode(&util.HttpResponse{
				Status:  http.StatusOK,
				Message: fmt.Sprintf("Registration token is: %s", viper.GetString("internaltoken")),
			})
		}
	})
}
