package notification

// Service is the interface definition for Notification sending services to allow abstraction from any particular
// implementation. As long as a service implements this interface it can send/receive messages to/from the client.
type Service interface {
	Receive()
	Send(Event)
	Hub() *Hub
}

type NotificationsService struct {
	hub       *Hub
	eventChan chan Event
	doneChan  chan struct{}
}

func NewNotificationsService(hub *Hub) *NotificationsService {
	ns := &NotificationsService{
		hub:       hub,
		eventChan: make(chan Event, 3),
		doneChan:  make(chan struct{}, 1),
	}

	go ns.hub.run()
	return ns
}

func (n *NotificationsService) Hub() *Hub {
	return n.hub
}

func (n *NotificationsService) Receive() {
}

func (n *NotificationsService) Send(event Event) {
	n.hub.broadcast <- event
}

func (n *NotificationsService) Close() {
	n.doneChan <- struct{}{}
}
