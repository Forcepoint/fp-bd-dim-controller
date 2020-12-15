package auth

import (
	"context"
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/auth/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strings"
)

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("x-access-token") //Grab the token from the header

		header = strings.TrimSpace(header)

		if header == "" {
			if r.URL.Path != "/api/ws" {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(util.HttpResponse{Status: http.StatusForbidden, Message: "Missing auth token"})
				return
			}
			token := r.URL.Query().Get("token")
			if token == "" {
				//Token is missing, returns with error code 403 Forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(util.HttpResponse{Status: http.StatusForbidden, Message: "Missing auth token"})
				return
			} else {
				header = strings.TrimSpace(token)
			}
		}

		tk := &structs.Token{}
		secret := os.Getenv("JWT_SECRET_KEY")

		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			//Token is expired, returns with error code 401 Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(util.HttpResponse{Status: http.StatusUnauthorized, Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func InternalAuthVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("x-internal-token") //Grab the token from the header

		header = strings.TrimSpace(header)

		if header == "" {
			//Token is missing, returns with error code 403 Forbidden
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(util.HttpResponse{Status: http.StatusForbidden, Message: "Missing auth token"})
			return
		}

		internalToken := viper.GetString("internaltoken")

		if header != internalToken {
			//Token is missing, returns with error code 401 Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(util.HttpResponse{Status: http.StatusUnauthorized, Message: "Unauthorized: Incorrect credentials"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
