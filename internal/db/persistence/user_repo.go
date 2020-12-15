package persistence

import (
	"database/sql"
	"fmt"
	"fp-dynamic-elements-manager-controller/internal/auth/structs"
	structs2 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	UserTable = "users"
)

type UserRepo struct {
	db  *sqlx.DB
	log *structs2.AppLogger
}

func NewUserRepo(appDb *sqlx.DB, logger *structs2.AppLogger) *UserRepo {
	return &UserRepo{db: appDb, log: logger}
}

func (u *UserRepo) InsertUser(item *structs.User) sql.Result {
	now := time.Now()

	smt := fmt.Sprintf("INSERT INTO %s (id, created_at, updated_at, deleted_at, name, email, password, admin) VALUES (?,?,?,?,?,?,?,?)", UserTable)
	tx, err := u.db.Begin()
	if err != nil {
		u.log.SystemLogger.Error(err, "Error starting transaction to insert user")
		return nil
	}
	res, err := tx.Exec(smt, item.ID, now, now, item.DeletedAt, item.Name, item.Email, item.Password, item.Admin)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			u.log.SystemLogger.Error(err, "Error inserting user, rolling back")
		}
		return nil
	}

	err = tx.Commit()

	if err != nil {
		u.log.SystemLogger.Error(err, "Error committing insert user")
		return nil
	}

	return res
}

func (u *UserRepo) UpdateUser(item *structs.User) sql.Result {
	var valueArgs []interface{}

	valueArgs = append(valueArgs, time.Now())
	valueArgs = append(valueArgs, item.Name)
	valueArgs = append(valueArgs, item.Email)
	valueArgs = append(valueArgs, item.Password)
	valueArgs = append(valueArgs, item.Email)

	smt := fmt.Sprintf("UPDATE %s SET updated_at = ?, name = ?, email = ?, password = ? WHERE email = ?", UserTable)
	tx, err := u.db.Begin()
	if err != nil {
		u.log.SystemLogger.Error(err, "Error starting transaction to insert user")
		return nil
	}
	res, err := tx.Exec(smt, valueArgs...)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			u.log.SystemLogger.Error(err, "Error inserting user, rolling back")
		}
		return nil
	}

	err = tx.Commit()

	if err != nil {
		u.log.SystemLogger.Error(err, "Error committing insert user")
		return nil
	}

	return res
}

func (u *UserRepo) DeleteByEmail(email string) error {
	smt := fmt.Sprintf(`DELETE FROM %s WHERE email = ?`, UserTable)
	tx, err := u.db.Begin()
	if err != nil {
		u.log.SystemLogger.Error(err, "Error starting transaction to delete list element")
		return err
	}
	_, err = tx.Exec(smt, email)
	if err != nil {
		u.log.SystemLogger.Error(err, "Error deleting list element, rolling back")
		err = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		u.log.SystemLogger.Error(err, "Error committing delete list element")
		return err
	}

	return nil
}

func (u *UserRepo) GetByEmail(email string) (receiver structs.User, err error) {
	err = u.db.Get(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE email = ? ORDER BY created_at DESC LIMIT 1;", UserTable), email)
	return
}

func (u *UserRepo) GetAll() (receiver []structs.ApiUser, err error) {
	err = u.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE deleted_at IS NULL AND admin = false ORDER BY created_at DESC;", UserTable))
	return
}

func (u *UserRepo) Exists(email string) bool {
	var usr structs.User
	return u.db.Get(&usr, fmt.Sprintf("SELECT * FROM %s WHERE email = ? LIMIT 1;", UserTable), email) == nil
}
