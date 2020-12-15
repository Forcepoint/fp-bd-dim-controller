package user

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	structs2 "fp-dynamic-elements-manager-controller/internal/auth/structs"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"fp-dynamic-elements-manager-controller/internal/notification"
	"fp-dynamic-elements-manager-controller/internal/user"
	"net/http"
)

// Handler handles all requests for the /user route
func Handler(repo *persistence.UserRepo, logger *structs.AppLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,POST,PUT,DELETE")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			users, err := user.GetAllUsers(repo)
			if err != nil {
				logger.SystemLogger.Error(err, "error getting all users")
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, "error retrieving resource")
				return
			}
			if users == nil {
				users = []structs2.ApiUser{}
			}
			json.NewEncoder(w).Encode(&util.HttpResponse{
				Items:   users,
				Status:  http.StatusOK,
				Message: "ok",
			})
		case http.MethodPost:
			usr := &structs2.User{}
			err := json.NewDecoder(r.Body).Decode(usr)
			if err != nil {
				util.ReturnHTTPStatus(w, http.StatusBadRequest, "bad request")
				return
			}
			userCreated := user.CreateUser(usr, repo, logger)
			if !userCreated {
				util.ReturnHTTPStatus(w, http.StatusOK, "user already exists")
				return
			}
			util.ReturnHTTPStatus(w, http.StatusCreated, "user created successfully")
		case http.MethodPut:
			usr := structs2.User{}
			err := json.NewDecoder(r.Body).Decode(&usr)
			if err != nil {
				logger.SystemLogger.Error(err, "error decoding json into entity")
				return
			}
			err = user.UpdateUserPassword(usr, repo, logger)
			if err != nil {
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, err.Error())
				return
			}
			util.ReturnHTTPStatus(w, http.StatusOK, "password changed successfully")
		case http.MethodDelete:
			usr := structs2.ApiUser{}
			err := json.NewDecoder(r.Body).Decode(&usr)
			if err != nil {
				util.ReturnHTTPStatus(w, http.StatusBadRequest, "bad request")
				return
			}
			err = user.DeleteUserByEmail(usr, repo, logger)
			if err != nil {
				util.ReturnHTTPStatus(w, http.StatusBadRequest, "bad request")
				return
			}
			logger.NotificationService.Send(notification.Event{
				EventType: notification.Success,
				Value:     "User deleted successfully",
			})
			util.ReturnHTTPStatus(w, http.StatusOK, "user deleted successfully")
		}
		return
	})
}
