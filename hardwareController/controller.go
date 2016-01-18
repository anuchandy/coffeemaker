package hardwareController

import (
	"github.com/anuchandy/coffeemaker/events"
	"github.com/anuchandy/coffeemaker/hardwareAPI"
	"github.com/anuchandy/coffeemaker/machineEvent"
	"log"
)

// Controller is the base type used by controllers specific to various
// components of Mark IV Coffee Maker.
type Controller struct {
	Api hardwareAPI.CommandAPI
}

// toMachineEvent returns coffee-machine event value stored in Event
// interface value e.
//
// It is panic if the event value is not MachineEvent.
func toMachineEvent(e events.Event) *machineEvent.MachineEvent {
	mEvent, ok := e.(*machineEvent.MachineEvent)
	if !ok {
		log.Panicf("Unknown event %v", e)
	}

	return mEvent
}
