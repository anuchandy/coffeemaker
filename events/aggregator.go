package events

// An Aggregator that implements pub-sub pattern.
type Aggregator struct {
	eventsChan  chan Event
	subscribers map[Event][]Subscriber
}

// NewAggregator returns an Aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{
		eventsChan:  make(chan Event),
		subscribers: make(map[Event][]Subscriber),
	}
}

// Publish publishes the event e, all subscribers subscribed for this will
// event will receive notification.
//
// Aggregator should be started by calling Start() method before trying to
// publish the first event, if the aggregator is not already started then
// calling go routine will be blocked.
//
// An attempt to publish an event after stopping the aggregator is a panic,
// aggregator is stopped by calling Stop() method.
func (a *Aggregator) Publish(e Event) {
	a.eventsChan <- e
}

// Subscribe register subscriber for the provided events es.
func (a *Aggregator) Subscribe(subscriber Subscriber, es ...Event) {
	for _, e := range es {
		a.subscribers[e] = append(a.subscribers[e], subscriber)
	}
}

// SubscribeFunc register the subscribe function subscriber for the provided list
// of events es.
func (a *Aggregator) SubscribeFunc(subscriber func(Event), es ...Event) {
	for _, e := range es {
		a.subscribers[e] = append(a.subscribers[e], SubscriberFunc(subscriber))
	}
}

// Start starts the aggregator, this enables aggregator to listen for published
// events and dispatches them to registered subscribers.
func (a *Aggregator) Start() {
	go func() {
		for {
			e, ok := <-a.eventsChan
			if !ok {
				return // channel closed by Stop method
			}

			// dispatch the event e
			for _, subscriber := range a.subscribers[e] {
				go subscriber.HandleEvent(e)
			}
		}
	}()
}

// Stop stops the aggregator, when the aggregator is in stopped state - a call to
// Publish() will be a panic and registered subscribers will not receive any more
// events.
func (a *Aggregator) Stop() {
	// Closing the event channel, any attempt to send to this channel by a publisher
	// will be panic.
	close(a.eventsChan)
}
