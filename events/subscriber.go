package events

// Subscriber is the interface that subscriber should satisfies
// inorder to be registered subscriber in the aggregator.
type Subscriber interface {
	HandleEvent(e Event)
}

type SubscriberFunc func(Event)

func (sf SubscriberFunc) HandleEvent(e Event) {
	sf(e)
}
