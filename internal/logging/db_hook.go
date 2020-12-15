package logging

import (
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/sirupsen/logrus"
)

// hook to buffer logs and only send at right severity.
type DatabaseHook struct {
	repo *persistence.LogEntryRepo
}

func NewDatabaseHook(repo *persistence.LogEntryRepo) *DatabaseHook {
	return &DatabaseHook{repo: repo}
}

// Fire will append all logs to a circular buffer and only 'flush'
// them when a log of sufficient severity(ERROR) is emitted.
func (h *DatabaseHook) Fire(entry *logrus.Entry) error {
	h.repo.InsertLogEntry(&structs.LogEntry{
		ModuleName: "Controller",
		Level:      entry.Level.String(),
		Message:    entry.Message,
		Caller:     entry.Caller.Func.Name(),
		Time:       entry.Time,
	})
	return nil
}

// Levels define on which log levels this DatabaseHook would trigger
func (h *DatabaseHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel}
}
