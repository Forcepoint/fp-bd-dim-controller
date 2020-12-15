package util

import (
	"encoding/json"
	"net/http"
)

func ReturnHTTPStatus(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&HttpResponse{
		Status:  status,
		Message: msg,
	})
}
