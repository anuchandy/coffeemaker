package hardwareController

import (
	"github.com/anuchandy/coffeemaker/events"
	"github.com/anuchandy/coffeemaker/hardwareAPI"
	"github.com/anuchandy/coffeemaker/machineEvent"
)

// WarmerPlateController is a controller for coffee-machine's warmer-plate.
type WarmerPlateController struct {
	Controller
}

// NewWarmerPlateController returns WarmerPlateController, the controller will be registered for the
// events those are used to set the coffee-maker's warmer-plate state (ON, OFF) using api.
func NewWarmerPlateController(aggregator *events.Aggregator, api hardwareAPI.CommandAPI) *WarmerPlateController {
	controller := WarmerPlateController{Controller{Api: api}}

	// WarmerPlateController is a subscriber for warmer-plate events
	aggregator.SubscribeFunc(controller.handleWarmerPlatePotNotEmptyEvent, machineEvent.WamerPlatePotNotEmpty)
	aggregator.SubscribeFunc(controller.handleEmptyEvents, machineEvent.WamerPlateEmpty, machineEvent.WamerPlatePotEmpty)

	return &controller
}

// handleWarmerPlatePotNotEmptyEvent is the subscribed event handler for pot-not-empty event.
func (c *WarmerPlateController) handleWarmerPlatePotNotEmptyEvent(e events.Event) {
	c.Api.SetWarmerPlateState(hardwareAPI.WARMER_ON)
}

// handleEmptyEvents is the subscribed event handler for warmer-plate empty events.
func (c *WarmerPlateController) handleEmptyEvents(e events.Event) {
	c.Api.SetWarmerPlateState(hardwareAPI.WARMER_OFF)
}
