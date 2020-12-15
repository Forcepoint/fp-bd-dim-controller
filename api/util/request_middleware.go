package util

import (
	"net/http"
	"os"
	"strconv"
)

// AddHeaders adds the content-type header for each request and the CORS header if in debug mode
func AddHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if getEnvDebugMode() {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("content-type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func getEnvDebugMode() bool {
	val := os.Getenv("DEBUG_MODE")
	if val == "" {
		return false
	}
	ret, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return ret
}
