package backup

import (
	"fmt"
	structs2 "fp-dynamic-elements-manager-controller/internal/backup/structs"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/docker"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Command string

const (
	Backup  Command = "backup"
	Restore Command = "restore"
)

type Provider interface {
	StartAutoBackup(Schedule) error
	Backup(string) error
	Restore(string) error
	List() ([]structs2.History, error)
}

type DatabaseBackupProvider struct {
	repo      persistence.ElementRepo
	docker    docker.Dockers
	logger    *structs.AppLogger
	committer HistoryCommitter
	scheduler *gocron.Scheduler
}

func NewDatabaseBackupProvider(
	docker docker.Dockers,
	committer HistoryCommitter,
	logger *structs.AppLogger,
	repo persistence.ElementRepo,
) *DatabaseBackupProvider {
	p := &DatabaseBackupProvider{
		repo:      repo,
		docker:    docker,
		logger:    logger,
		committer: committer,
		scheduler: gocron.NewScheduler(time.UTC),
	}

	return p
}

func (d *DatabaseBackupProvider) StartAutoBackup(schedule Schedule) error {
	d.scheduler.Clear()
	w, err := parseWeekday(strings.ToLower(schedule.DayOfWeek))
	if err != nil {
		return err
	}
	d.scheduler.Every(1).Weekday(w).At(schedule.TimeOfDay).Do(d.Backup, "Auto")
	// scheduler starts running jobs and current thread continues to execute
	d.scheduler.StartAsync()
	writeSchedule(schedule)
	return nil
}

func (d *DatabaseBackupProvider) Backup(message string) error {
	err := d.docker.RunDatabaseDump()
	if err != nil {
		return err
	}
	count, err := d.repo.GetTotalElementCount()
	if err != nil {
		return err
	}
	err = d.committer.Commit(message, count)
	if err != nil {
		return err
	}
	return nil
}

func (d *DatabaseBackupProvider) Restore(commitHash string) error {
	err := d.committer.RestoreToPoint(commitHash)
	if err != nil {
		d.logger.SystemLogger.Error(err, fmt.Sprintf("error rolling back to commit: %s", commitHash))
		return err
	}
	d.docker.RunDatabaseRestore()
	return nil
}

func (d *DatabaseBackupProvider) List() ([]structs2.History, error) {
	return d.committer.ListHistory()
}

func writeSchedule(schedule Schedule) {
	viper.Set("dayofweek", schedule.DayOfWeek)
	viper.Set("timeofday", schedule.TimeOfDay)

	if err := viper.WriteConfig(); err != nil {
		log.Error().Err(err).Msg("error writing schedule")
	}
}

type Schedule struct {
	DayOfWeek string `json:"day"`
	TimeOfDay string `json:"time"`
}
