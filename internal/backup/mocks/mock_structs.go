package mocks

import (
	"fp-dynamic-elements-manager-controller/internal/backup/structs"
	"fp-dynamic-elements-manager-controller/internal/notification"
	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Panic(s string) {
	l.Called(s)
}

func (l *LoggerMock) Fatal(err error, s string) {
	l.Called(err, s)
}

func (l *LoggerMock) Error(err error, s string) {
	l.Called(err, s)
}

func (l *LoggerMock) Warn(s string) {
	l.Called(s)
}

func (l *LoggerMock) Info(s string) {
	l.Called(s)
}

func (l *LoggerMock) Debug(s string) {
	l.Called(s)
}

func (l *LoggerMock) Trace(s string) {
	l.Called(s)
}

type RepoMock struct {
	mock.Mock
}

func (r *RepoMock) GetTotalElementCount() (int64, error) {
	args := r.Called()
	return int64(args.Int(0)), args.Error(1)
}

type DockerMock struct {
	mock.Mock
}

func (d *DockerMock) PullAndStart(a, b string) error {
	args := d.Called(a, b)
	return args.Error(0)
}

func (d *DockerMock) Create(a, b, c string, arr1 []string, arr2 []string) error {
	args := d.Called(a, b, c, arr1, arr2)
	return args.Error(1)
}

func (d *DockerMock) Start(s string) error {
	args := d.Called(s)
	return args.Error(0)
}

func (d *DockerMock) Stop(s string) error {
	args := d.Called(s)
	return args.Error(0)
}

func (d *DockerMock) Restart(s string) error {
	args := d.Called(s)
	return args.Error(0)
}

func (d *DockerMock) Remove(s string) error {
	args := d.Called(s)
	return args.Error(0)
}

func (d *DockerMock) ListNetworks() []types.NetworkResource {
	args := d.Called()
	return args.Get(0).([]types.NetworkResource)
}

func (d *DockerMock) ListContainers(options types.ContainerListOptions) []types.Container {
	args := d.Called(options)
	return args.Get(0).([]types.Container)
}

func (d *DockerMock) RunDatabaseDump() error {
	args := d.Called()
	return args.Error(0)
}

func (d *DockerMock) RunDatabaseRestore() {
	d.Called()
}

type NSMock struct {
	mock.Mock
}

func (n *NSMock) Receive() {
	n.Called()
}

func (n *NSMock) Send(event notification.Event) {
	n.Called()
}

func (n *NSMock) Hub() *notification.Hub {
	n.Called()
	return nil
}

type CommitterMock struct {
	mock.Mock
}

func (c *CommitterMock) Commit(s string, i int64) error {
	args := c.Called(s, i)
	return args.Error(0)
}

func (c *CommitterMock) RestoreToPoint(s string) error {
	args := c.Called(s)
	return args.Error(0)
}

func (c *CommitterMock) ListHistory() ([]structs.History, error) {
	args := c.Called()
	return args.Get(0).([]structs.History), args.Error(1)
}
