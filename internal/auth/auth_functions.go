package auth

import (
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/auth/structs"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	validation "fp-dynamic-elements-manager-controller/internal/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func FindUserAndCreateToken(email, password string, repo *persistence.UserRepo) (*util.HttpResponse, map[string]interface{}) {
	if !validation.IsEmailValid(email) {
		return &util.HttpResponse{Status: http.StatusUnauthorized, Message: "Invalid email format."}, nil
	}
	user, err := repo.GetByEmail(email)

	if err != nil {
		return &util.HttpResponse{Status: http.StatusUnauthorized, Message: "Invalid login credentials."}, nil
	}
	expiresAt := time.Now().Add(time.Hour * 12).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return &util.HttpResponse{Status: http.StatusUnauthorized, Message: "Invalid login credentials."}, nil
	}

	tk := &structs.Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)

	secret := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		logrus.Error(err)
	}

	var resp = map[string]interface{}{"status": true, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = map[string]interface{}{"name": user.Name, "email": user.Email}
	return nil, resp
}
