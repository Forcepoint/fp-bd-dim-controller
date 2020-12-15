package docker

import (
	"fmt"
	"fp-dynamic-elements-manager-controller/internal/docker/mocks"
	"fp-dynamic-elements-manager-controller/internal/docker/structs"
	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var testContainer = structs.ContainerDetails{
	ID:                "fp-test",
	Name:              "FP Testing Module",
	Description:       "This is just to run tests",
	Type:              structs.INGRESS,
	AcceptedDataTypes: []structs.ElementType{structs.IP, structs.RANGE},
	Volumes:           []string{"/config:/config"},
	Network:           "module_net",
	EnvVars:           []string{"LOCAL_PORT=8080"},
	ImageRef:          "test.docker.io/fp-test/test:latest",
	IconURL:           "www.icon.org",
	Command:           structs.PullAndStart,
	RegistrationToken: "123456",
}

type CommandHandlerTestSuite struct {
	suite.Suite
}

func TestCommandHandler(t *testing.T) {
	suite.Run(t, new(CommandHandlerTestSuite))
}

func (c *CommandHandlerTestSuite) TestCommandHandler_Create() {
	modRepo := new(mocks.MockModuleMetadataRepo)
	docker := new(mocks.TestDocker)

	testContainer.Command = structs.Create

	docker.On("ListNetworks").Return([]types.NetworkResource{{
		Name: "module_net",
		ID:   "module_net",
	}})

	docker.On(
		"Create",
		testContainer.ImageRef,
		testContainer.ID,
		testContainer.Network,
		testContainer.Volumes,
		mock.Anything,
	).Return(nil)

	handler := NewCommandHandler(docker, modRepo)

	doneCh, evtCh, errCh := handler.RunCommands(structs.ContainerDetailsWrapper{
		Containers: []structs.ContainerDetails{
			testContainer,
		}})

readChannel:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				c.T().Error("test failed", err)
			}
		case event := <-evtCh:
			fmt.Println(event)
		case <-doneCh:
			break readChannel
		}
	}

	docker.AssertCalled(
		c.T(),
		"Create",
		testContainer.ImageRef,
		testContainer.ID,
		testContainer.Network,
		testContainer.Volumes,
		mock.Anything,
	)

	docker.AssertExpectations(c.T())
}

func (c *CommandHandlerTestSuite) TestCommandHandler_PullAndStart() {
	modRepo := new(mocks.MockModuleMetadataRepo)
	docker := new(mocks.TestDocker)

	testContainer.Command = structs.PullAndStart

	docker.On("ListNetworks").Return([]types.NetworkResource{{
		Name: "module_net",
		ID:   "module_net",
	}})

	docker.On("PullAndStart", testContainer.ImageRef, testContainer.ID).Return(nil)

	handler := NewCommandHandler(docker, modRepo)

	doneCh, evtCh, errCh := handler.RunCommands(structs.ContainerDetailsWrapper{
		Containers: []structs.ContainerDetails{
			testContainer,
		}})

readChannel:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				c.T().Error("test failed", err)
			}
		case event := <-evtCh:
			fmt.Println(event)
		case <-doneCh:
			break readChannel
		}
	}

	docker.AssertCalled(c.T(), "PullAndStart", testContainer.ImageRef, testContainer.ID)

	docker.AssertExpectations(c.T())
}

func (c *CommandHandlerTestSuite) TestCommandHandler_Start() {
	modRepo := new(mocks.MockModuleMetadataRepo)
	docker := new(mocks.TestDocker)

	testContainer.Command = structs.Start

	docker.On("ListNetworks").Return([]types.NetworkResource{{
		Name: "module_net",
		ID:   "module_net",
	}})

	docker.On("Start", testContainer.ID).Return(nil)

	handler := NewCommandHandler(docker, modRepo)

	doneCh, evtCh, errCh := handler.RunCommands(structs.ContainerDetailsWrapper{
		Containers: []structs.ContainerDetails{
			testContainer,
		}})

readChannel:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				c.T().Error("test failed", err)
			}
		case event := <-evtCh:
			fmt.Println(event)
		case <-doneCh:
			break readChannel
		}
	}

	docker.AssertCalled(c.T(), "Start", testContainer.ID)

	docker.AssertExpectations(c.T())
}

func (c *CommandHandlerTestSuite) TestCommandHandler_Stop() {
	modRepo := new(mocks.MockModuleMetadataRepo)
	docker := new(mocks.TestDocker)

	testContainer.Command = structs.Stop

	docker.On("ListNetworks").Return([]types.NetworkResource{{
		Name: "module_net",
		ID:   "module_net",
	}})

	docker.On("Stop", testContainer.ID).Return(nil)

	handler := NewCommandHandler(docker, modRepo)

	doneCh, evtCh, errCh := handler.RunCommands(structs.ContainerDetailsWrapper{
		Containers: []structs.ContainerDetails{
			testContainer,
		}})

readChannel:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				c.T().Error("test failed", err)
			}
		case event := <-evtCh:
			fmt.Println(event)
		case <-doneCh:
			break readChannel
		}
	}

	docker.AssertCalled(c.T(), "Stop", testContainer.ID)

	docker.AssertExpectations(c.T())
}

func (c *CommandHandlerTestSuite) TestCommandHandler_Restart() {
	modRepo := new(mocks.MockModuleMetadataRepo)
	docker := new(mocks.TestDocker)

	testContainer.Command = structs.Restart

	docker.On("ListNetworks").Return([]types.NetworkResource{{
		Name: "module_net",
		ID:   "module_net",
	}})

	docker.On("Restart", testContainer.ID).Return(nil)

	handler := NewCommandHandler(docker, modRepo)

	doneCh, evtCh, errCh := handler.RunCommands(structs.ContainerDetailsWrapper{
		Containers: []structs.ContainerDetails{
			testContainer,
		}})

readChannel:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				c.T().Error("test failed", err)
			}
		case event := <-evtCh:
			fmt.Println(event)
		case <-doneCh:
			break readChannel
		}
	}

	docker.AssertCalled(c.T(), "Restart", testContainer.ID)

	docker.AssertExpectations(c.T())
}

func (c *CommandHandlerTestSuite) TestCommandHandler_Remove() {
	modRepo := new(mocks.MockModuleMetadataRepo)
	docker := new(mocks.TestDocker)

	testContainer.Command = structs.Remove

	docker.On("ListNetworks").Return([]types.NetworkResource{{
		Name: "module_net",
		ID:   "module_net",
	}})

	docker.On("Remove", testContainer.ID).Return(nil)

	modRepo.On("DeleteByServiceName", mock.Anything).Return(nil)

	handler := NewCommandHandler(docker, modRepo)

	doneCh, evtCh, errCh := handler.RunCommands(structs.ContainerDetailsWrapper{
		Containers: []structs.ContainerDetails{
			testContainer,
		}})

readChannel:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				c.T().Error("test failed", err)
			}
		case event := <-evtCh:
			fmt.Println(event)
		case <-doneCh:
			break readChannel
		}
	}

	docker.AssertCalled(c.T(), "Remove", testContainer.ID)
	modRepo.AssertCalled(c.T(), "DeleteByServiceName", mock.Anything)

	docker.AssertExpectations(c.T())
	modRepo.AssertExpectations(c.T())
}

func (c *CommandHandlerTestSuite) TestCommandHandler_MultipleCommands() {
	modRepo := new(mocks.MockModuleMetadataRepo)
	docker := new(mocks.TestDocker)

	testContainer.Command = structs.Create
	testContainer2 := testContainer
	testContainer.Command = structs.PullAndStart
	testContainer3 := testContainer
	testContainer.Command = structs.Start
	testContainer4 := testContainer
	testContainer.Command = structs.Stop
	testContainer5 := testContainer
	testContainer.Command = structs.Restart
	testContainer6 := testContainer
	testContainer.Command = structs.Remove

	docker.On("ListNetworks").Times(6).Return([]types.NetworkResource{{
		Name: "module_net",
		ID:   "module_net",
	}})

	docker.On(
		"Create",
		testContainer.ImageRef,
		testContainer.ID,
		testContainer.Network,
		testContainer.Volumes,
		mock.Anything,
	).Times(1).Return(nil)

	docker.On("PullAndStart", testContainer.ImageRef, testContainer.ID).Times(1).Return(nil)
	docker.On("Start", testContainer.ID).Times(1).Return(nil)
	docker.On("Stop", testContainer.ID).Times(1).Return(nil)
	docker.On("Restart", testContainer.ID).Times(1).Return(nil)
	docker.On("Remove", testContainer.ID).Times(1).Return(nil)

	modRepo.On("DeleteByServiceName", mock.Anything).Times(1).Return(nil)

	handler := NewCommandHandler(docker, modRepo)

	doneCh, evtCh, errCh := handler.RunCommands(structs.ContainerDetailsWrapper{
		Containers: []structs.ContainerDetails{
			testContainer,
			testContainer2,
			testContainer3,
			testContainer4,
			testContainer5,
			testContainer6,
		}})

readChannel:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				c.T().Error("test failed", err)
			}
		case event := <-evtCh:
			fmt.Println(event)
		case <-doneCh:
			break readChannel
		}
	}

	docker.AssertCalled(
		c.T(),
		"Create",
		testContainer.ImageRef,
		testContainer.ID,
		testContainer.Network,
		testContainer.Volumes,
		mock.Anything,
	)

	docker.AssertCalled(c.T(), "PullAndStart", testContainer.ImageRef, testContainer.ID)

	docker.AssertCalled(c.T(), "Start", testContainer.ID)

	docker.AssertCalled(c.T(), "Stop", testContainer.ID)

	docker.AssertCalled(c.T(), "Restart", testContainer.ID)

	docker.AssertCalled(c.T(), "Remove", testContainer.ID)
	modRepo.AssertCalled(c.T(), "DeleteByServiceName", mock.Anything)

	docker.AssertExpectations(c.T())
	modRepo.AssertExpectations(c.T())
}
