package structs

import "fp-dynamic-elements-manager-controller/internal/notification"

type ModuleType string
type ElementType string
type Command string

const (
	PullAndStart   Command = "pullstart"
	PullAndRestart Command = "pullrestart"
	Start          Command = "start"
	Stop           Command = "stop"
	Restart        Command = "restart"
	Create         Command = "create"
	Remove         Command = "remove"
	List           Command = "list"

	INGRESS    ModuleType = "ingress"
	EGRESS     ModuleType = "egress"
	FUNCTIONAL ModuleType = "functional"

	IP     ElementType = "IP"
	DOMAIN ElementType = "DOMAIN"
	URL    ElementType = "URL"
	RANGE  ElementType = "RANGE"
)

type ContainerDetails struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	Type              ModuleType    `json:"type"`
	AcceptedDataTypes []ElementType `json:"accepted_data_types"`
	Volumes           []string      `json:"volumes"`
	Network           string        `json:"network"`
	EnvVars           []string      `json:"env_vars"`
	ImageRef          string        `json:"image_ref"`
	IconURL           string        `json:"icon_url"`
	Command           Command       `json:"command"`
	RegistrationToken string        `json:"registration_token"`
}

type ContainerDetailsWrapper struct {
	Containers []ContainerDetails `json:"containers"`
}

type ContainerVolume struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type ContainerNames struct {
	Containers []string `json:"containers"`
}

func (c Command) CommandToState() notification.State {
	switch c {
	case Create, PullAndStart:
		return notification.Created
	case Stop:
		return notification.Stopped
	case Start, Restart:
		return notification.Started
	case Remove:
		return notification.Deleted
	}
	return notification.None
}
