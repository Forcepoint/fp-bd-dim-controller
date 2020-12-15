package structs

import (
	notificationfuncs "fp-dynamic-elements-manager-controller/internal/notification"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"os"
)

type Logger interface {
	Panic(string)
	Fatal(error, string)
	Error(error, string)
	Warn(string)
	Info(string)
	Debug(string)
	Trace(string)
}

type AppLogger struct {
	UserLogger          Logger
	SystemLogger        Logger
	NotificationService notificationfuncs.Service
}

func NewAppLogger(service notificationfuncs.Service) *AppLogger {
	return &AppLogger{
		UserLogger:          &UserLogger{log: logrus.StandardLogger()},
		SystemLogger:        &SystemLogger{log: zerolog.New(os.Stderr)},
		NotificationService: service,
	}
}

// Define system logger, these logs do not get pushed to the log table,
//these are for developer use
type SystemLogger struct {
	log zerolog.Logger
}

func (s *SystemLogger) Panic(msg string) {
	s.log.Panic().Msg(msg)
}

func (s *SystemLogger) Fatal(err error, msg string) {
	s.log.Fatal().Err(err).Msg(msg)
}

func (s *SystemLogger) Error(err error, msg string) {
	s.log.Error().Err(err).Msg(msg)
}

func (s *SystemLogger) Warn(msg string) {
	s.log.Warn().Msg(msg)
}

func (s *SystemLogger) Info(msg string) {
	s.log.Info().Msg(msg)
}

func (s *SystemLogger) Debug(msg string) {
	s.log.Debug().Msg(msg)
}

func (s *SystemLogger) Trace(msg string) {
	s.log.Trace().Msg(msg)
}

// Define user logger, these logs do get pushed to the log table,
//these are for useful messages for the end user
type UserLogger struct {
	log *logrus.Logger
}

func (u *UserLogger) Panic(msg string) {
	u.log.Panic(msg)
}

func (u *UserLogger) Fatal(err error, msg string) {
	u.log.Fatal(err, msg)
}

func (u *UserLogger) Error(err error, msg string) {
	u.log.Error(err, msg)
}

func (u *UserLogger) Warn(msg string) {
	u.log.Warn(msg)
}

func (u *UserLogger) Info(msg string) {
	u.log.Info(msg)
}

func (u *UserLogger) Debug(msg string) {
	u.log.Debug(msg)
}

func (u *UserLogger) Trace(msg string) {
	u.log.Trace(msg)
}
