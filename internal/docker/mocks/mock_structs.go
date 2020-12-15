package mocks

import (
	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/mock"
)

type MockModuleMetadataRepo struct {
	mock.Mock
}

func (r *MockModuleMetadataRepo) DeleteByServiceName(svcName string) error {
	args := r.Called(svcName)
	return args.Error(0)
}

type TestDocker struct {
	mock.Mock
}

func (t *TestDocker) PullAndStart(ref, id string) error {
	args := t.Called(ref, id)
	return args.Error(0)
}

func (t *TestDocker) Create(imageRef, containerName, containerNetwork string, volumes []string, envVars []string) error {
	args := t.Called(imageRef, containerName, containerNetwork, volumes, envVars)
	return args.Error(0)
}

func (t *TestDocker) Start(id string) error {
	args := t.Called(id)
	return args.Error(0)
}

func (t *TestDocker) Stop(id string) error {
	args := t.Called(id)
	return args.Error(0)
}

func (t *TestDocker) Restart(id string) error {
	args := t.Called(id)
	return args.Error(0)
}

func (t *TestDocker) Remove(id string) error {
	args := t.Called(id)
	return args.Error(0)
}

func (t *TestDocker) ListNetworks() []types.NetworkResource {
	args := t.Called()
	return args.Get(0).([]types.NetworkResource)
}

func (t *TestDocker) ListContainers(opts types.ContainerListOptions) []types.Container {
	args := t.Called(opts)
	return args.Get(0).([]types.Container)
}

func (t *TestDocker) RunDatabaseDump() error {
	args := t.Called()
	return args.Error(0)
}

func (t *TestDocker) RunDatabaseRestore() {
	t.Called()
}
