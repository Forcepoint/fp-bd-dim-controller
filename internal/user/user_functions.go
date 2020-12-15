package user

import (
	"errors"
	"fp-dynamic-elements-manager-controller/internal/auth/structs"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	structs2 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	notificationfuncs "fp-dynamic-elements-manager-controller/internal/notification"
	"fp-dynamic-elements-manager-controller/internal/util"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidEmailFormat = errors.New("invalid email format")

func CreateUser(user *structs.User, userRepo *persistence.UserRepo, logger *structs2.AppLogger) bool {
	if !util.IsEmailValid(user.Email) {
		logger.NotificationService.Send(notificationfuncs.Event{
			EventType: notificationfuncs.Error,
			Value:     "Invalid email format",
		})
		return false
	}
	if userRepo.Exists(user.Email) {
		logger.NotificationService.Send(notificationfuncs.Event{
			EventType: notificationfuncs.Info,
			Value:     "User already exists",
		})
		return false
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.SystemLogger.Error(err, "error encoding password for insert bcrypt")
		return false
	}
	user.Password = string(pass)
	result := userRepo.InsertUser(user)
	if result == nil {
		return false
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false
	}
	return rows != 0
}

func GetAllUsers(userRepo *persistence.UserRepo) ([]structs.ApiUser, error) {
	return userRepo.GetAll()
}

func UpdateUserPassword(user structs.User, userRepo *persistence.UserRepo, logger *structs2.AppLogger) error {
	if !util.IsEmailValid(user.Email) {
		logger.NotificationService.Send(notificationfuncs.Event{
			EventType: notificationfuncs.Error,
			Value:     "Invalid email format",
		})
		return ErrInvalidEmailFormat
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error encoding password for update bcrypt")
		return err
	}
	dbUser, err := userRepo.GetByEmail(user.Email)
	if err != nil {
		return err
	}
	dbUser.Password = string(pass)
	result := userRepo.UpdateUser(&dbUser)
	if result == nil {
		return errors.New("result from update was nil")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("result from update was 0 rows updated")
	}
	return nil
}

func CreateAdminUserIfNotExists(userRepo *persistence.UserRepo) error {
	user, err := userRepo.GetByEmail("admin.user@forcepoint.com")

	if err != nil || user.ID == 0 {
		// TODO change call above to an 'exists' call, extract hardcoded values here to the docker-compose file, change logger
		logrus.Info("Creating new Admin User")
		user := &structs.User{Name: "Admin User", Email: "admin.user@forcepoint.com", Password: "password1", Admin: true}

		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			panic(err)
		}

		user.Password = string(pass)
		result := userRepo.InsertUser(user)

		if result == nil {
			return errors.New("result from update was nil")
		}

		rows, err := result.RowsAffected()

		if err != nil {
			return err
		}

		if rows == 0 {
			return errors.New("result from update was 0 rows updated")
		}

		return nil
	}

	return nil
}

func DeleteUserByEmail(user structs.ApiUser, repo *persistence.UserRepo, logger *structs2.AppLogger) error {
	if !util.IsEmailValid(user.Email) {
		logger.NotificationService.Send(notificationfuncs.Event{
			EventType: notificationfuncs.Error,
			Value:     "Invalid email format",
		})
		return ErrInvalidEmailFormat
	}
	return repo.DeleteByEmail(user.Email)
}
