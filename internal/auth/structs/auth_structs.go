package structs

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	Name      string
	Email     string `json:"email"`
	Password  string `json:"password"`
	Admin     bool   `json:"admin"`
}

type ApiUser struct {
	ID        uint       `json:"-"`
	CreatedAt time.Time  `json:"-" db:"created_at"`
	UpdatedAt time.Time  `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Admin     bool       `json:"admin"`
}

type Token struct {
	UserID              uint   `json:"user_id"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	*jwt.StandardClaims `json:"standard_claims"`
}
