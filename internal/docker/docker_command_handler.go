package docker

import (
	"container/list"
	"fmt"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/docker/structs"
	"fp-dynamic-elements-manager-controller/internal/docker/utils"
	"fp-dynamic-elements-manager-controller/internal/notification"
	"github.com/docker/docker/api/types"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
)

type CommandHandler struct {
	mu           *sync.Mutex
	docker       Dockers
	commandQueue *list.List
	errCh        chan error
	doneCh       chan struct{}
	evtCh        chan notification.Event
	repo         persistence.ModuleRepo
}

func NewCommandHandler(d Dockers, mRepo persistence.ModuleRepo) *CommandHandler {
	return &CommandHandler{
		docker:       d,
		commandQueue: list.New(),
		mu:           &sync.Mutex{},
		errCh:        make(chan error),
		doneCh:       make(chan struct{}),
		evtCh:        make(chan notification.Event),
		repo:         mRepo,
	}
}

func (c *CommandHandler) RunCommands(details structs.ContainerDetailsWrapper) (<-chan struct{}, <-chan notification.Event, <-chan error) {
	// Add commands to the queue
	for _, container := range details.Containers {

		if container.Command == structs.Create {
			// Make sure the image that is trying to be used is coming from our own registry
			if !utils.ValidImageRef(container.ImageRef) {
				continue
			}
		}

		if err := enrichContainerStruct(&container, c.docker); err != nil {
			log.Error().Err(err)
			continue
		}

		c.commandQueue.PushBack(container)
	}

	// Start processing the commands
	go c.run()

	return c.doneCh, c.evtCh, c.errCh
}

func (c *CommandHandler) List() []types.Container {
	return c.docker.ListContainers(types.ContainerListOptions{})
}

func (c *CommandHandler) MapContainerNames() structs.ContainerNames {
	containers := c.List()

	var names []string
	for _, con := range containers {
		names = append(names, con.Names[0])
	}

	return structs.ContainerNames{Containers: names}
}

func (c *CommandHandler) run() {
	if c.commandQueue.Len() < 1 {
		c.doneCh <- struct{}{}
		return
	}

	// Create a lock here to avoid weird race conditions if new commands are pushed before old ones are done
	c.mu.Lock()
	defer c.mu.Unlock()

	// Pull from the queue and process the command
	for c.commandQueue.Len() > 0 {
		element := c.commandQueue.Front()
		casted := element.Value.(structs.ContainerDetails)
		err := c.processCommand(casted)
		c.commandQueue.Remove(element)

		if err == nil {
			c.evtCh <- notification.Event{
				EventType: notification.Success,
				Value:     fmt.Sprintf("%s module %s", strings.Title(string(casted.Command.CommandToState())), casted.Name),
				Context: notification.EventContext{
					Type:       notification.Module,
					Identifier: casted.ID,
					State:      casted.Command.CommandToState(),
				},
			}
		}
	}

	c.doneCh <- struct{}{}
}

func (c *CommandHandler) processCommand(container structs.ContainerDetails) error {
	var err error
	switch container.Command {
	case structs.PullAndStart:
		err = c.docker.PullAndStart(container.ImageRef, container.ID)
	case structs.Create:
		err = c.docker.Create(container.ImageRef, container.ID, container.Network, container.Volumes, container.EnvVars)
	case structs.Stop:
		err = c.docker.Stop(container.ID)
	case structs.Start:
		err = c.docker.Start(container.ID)
	case structs.Restart:
		err = c.docker.Restart(container.ID)
	case structs.Remove:
		err = c.docker.Remove(container.ID)
		if err == nil {
			err = c.deleteFromControllerDB(container.ID)
		}
	}

	if err != nil {
		c.errCh <- err
	}

	return err
}

func (c *CommandHandler) deleteFromControllerDB(svcName string) error {
	return c.repo.DeleteByServiceName(svcName)
}

func enrichContainerStruct(container *structs.ContainerDetails, docker Dockers) error {
	utils.BuildModuleEnvVars(&container.EnvVars)
	utils.AddModuleBindPaths(&container.Volumes, container.ID)

	network, err := utils.AddModuleNetwork(docker.ListNetworks())

	if err != nil {
		return err
	}

	container.Network = network

	return nil
}
