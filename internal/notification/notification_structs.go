package notification

type EventType string
type EntityType string
type State string

const (
	Log     EventType = "log"
	Error   EventType = "error"
	Warning EventType = "warning"
	Success EventType = "success"
	Info    EventType = "info"

	Module      EntityType = "module"
	ListElement EntityType = "listElement"

	None    State = "none"
	Deleted State = "deleted"
	Created State = "created"
	Started State = "started"
	Stopped State = "stopped"
)

type Event struct {
	EventType EventType    `json:"event_type"`
	Value     string       `json:"value"`
	Context   EventContext `json:"context"`
}

type EventContext struct {
	Type       EntityType `json:"type"`
	Identifier string     `json:"identifier"`
	State      State      `json:"state"`
}
