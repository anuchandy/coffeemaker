package hardwareController

import (
	"github.com/anuchandy/coffeemaker/events"
	"github.com/anuchandy/coffeemaker/hardwareAPI"
	"github.com/anuchandy/coffeemaker/machineEvent"
)

// BoilerController is a controller for coffee-machine's boiler heater.
type BoilerController struct {
	Controller
	boilerEmpty bool
}

// NewBoilerController returns BoilerController, the controller will be registered for the
// events those are used to set the coffee-maker's boiler heater state (ON, OFF) using api.
func NewBoilerController(aggregator *events.Aggregator, api hardwareAPI.CommandAPI) *BoilerController {
	controller := BoilerController{
		Controller:  Controller{Api: api},
		boilerEmpty: false,
	}

	// BoilerController is a subscriber for boiler events
	aggregator.SubscribeFunc(controller.handleBoilerEvents, machineEvent.BoilerEmpty, machineEvent.BoilerNotEmpty)
	// BoilerController is a subscriber for brew button event BrewButtonPushed
	aggregator.SubscribeFunc(controller.handleBrewButtonPushEvent, machineEvent.BrewButtonPushed)

	return &controller
}

// handleBoilerEvents is the subscribed event handler for boiler events.
func (c *BoilerController) handleBoilerEvents(e events.Event) {
	c.boilerEmpty = toMachineEvent(e) == machineEvent.BoilerEmpty
	if c.boilerEmpty {
		c.Api.SetBoilerState(hardwareAPI.BOILER_OFF)
	}
}

// handleBrewButtonPushEvent is the subscribed event handler for brew button pushed event.
func (c *BoilerController) handleBrewButtonPushEvent(e events.Event) {
	if !c.boilerEmpty {
		c.Api.SetBoilerState(hardwareAPI.BOILER_ON)
	}
}
