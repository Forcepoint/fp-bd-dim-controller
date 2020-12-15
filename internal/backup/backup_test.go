package backup

import (
	"errors"
	"fp-dynamic-elements-manager-controller/internal/backup/mocks"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BackupRestoreTestSuite struct {
	suite.Suite
}

func (b *BackupRestoreTestSuite) TestDatabaseBackupProvider_Backup() {
	loggerObj := new(mocks.LoggerMock)
	nsObj := new(mocks.NSMock)
	logger := &structs.AppLogger{
		UserLogger:          loggerObj,
		SystemLogger:        loggerObj,
		NotificationService: nsObj,
	}
	b.T().Run("Test Database Backup Provider Run (Backup) - No Errors", func(t *testing.T) {
		dockerObj := new(mocks.DockerMock)
		repoObj := new(mocks.RepoMock)
		committerObj := new(mocks.CommitterMock)
		// setup expectations
		dockerObj.On("RunDatabaseDump").Times(1).Return(nil)
		repoObj.On("GetTotalElementCount").Times(1).Return(5, nil)
		committerObj.On("Commit", int64(5)).Times(1).Return(nil)

		provider := NewDatabaseBackupProvider(dockerObj, committerObj, logger, repoObj)

		assert.Nil(b.T(), provider.Backup("Manual"))

		dockerObj.AssertCalled(b.T(), "RunDatabaseDump")
		repoObj.AssertCalled(b.T(), "GetTotalElementCount")
		committerObj.AssertCalled(b.T(), "Commit", int64(5))

		// assert that the expectations were met
		repoObj.AssertExpectations(b.T())
		repoObj.AssertExpectations(b.T())
		committerObj.AssertExpectations(b.T())
	})

	b.T().Run("Test Database Backup Provider Run (Backup) - Docker Error", func(t *testing.T) {
		dockerObj := new(mocks.DockerMock)
		repoObj := new(mocks.RepoMock)
		committerObj := new(mocks.CommitterMock)
		// setup expectations
		dockerObj.On("RunDatabaseDump").Times(1).Return(errors.New("docker error"))

		provider := NewDatabaseBackupProvider(dockerObj, committerObj, logger, repoObj)

		assert.NotNil(b.T(), provider.Backup("Manual"))

		repoObj.AssertNotCalled(b.T(), "GetTotalElementCount", int64(5))
		committerObj.AssertNotCalled(b.T(), "Commit")

		// assert that the expectations were met
		repoObj.AssertExpectations(b.T())
		committerObj.AssertExpectations(b.T())
	})

	b.T().Run("Test Database Backup Provider (Backup) - Commit Error", func(t *testing.T) {
		dockerObj := new(mocks.DockerMock)
		repoObj := new(mocks.RepoMock)
		committerObj := new(mocks.CommitterMock)
		// setup expectations
		dockerObj.On("RunDatabaseDump").Times(1).Return(nil)
		repoObj.On("GetTotalElementCount").Times(1).Return(5, nil)
		committerObj.On("Commit", int64(5)).Times(1).Return(errors.New("commit error"))

		provider := NewDatabaseBackupProvider(dockerObj, committerObj, logger, repoObj)

		assert.NotNil(b.T(), provider.Backup("Manual"))

		committerObj.AssertNotCalled(b.T(), "Commit")

		// assert that the expectations were met
		committerObj.AssertExpectations(b.T())
	})

	b.T().Run("Test Database Backup Provider (Backup) - Repo Error", func(t *testing.T) {
		dockerObj := new(mocks.DockerMock)
		repoObj := new(mocks.RepoMock)
		committerObj := new(mocks.CommitterMock)

		// setup expectations
		dockerObj.On("RunDatabaseDump").Times(1).Return(nil)
		repoObj.On("GetTotalElementCount").Times(1).Return(5, errors.New("repo error"))

		provider := NewDatabaseBackupProvider(dockerObj, committerObj, logger, repoObj)

		assert.NotNil(b.T(), provider.Backup("Manual"))

		repoObj.AssertCalled(b.T(), "GetTotalElementCount")
		committerObj.AssertNotCalled(b.T(), "Commit", int64(5))

		// assert that the expectations were met
		committerObj.AssertExpectations(b.T())
	})
}

func TestDatabaseBackupProvider(t *testing.T) {
	suite.Run(t, new(BackupRestoreTestSuite))
}
