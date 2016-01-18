package events

// Event is the interface that an event type should satisfies
// inorder to be used as an event by publisher.
type Event interface {
	// Event returns string representation of the event.
	Event() string
}
