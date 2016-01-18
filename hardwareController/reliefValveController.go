package hardwareController

import (
	"github.com/anuchandy/coffeemaker/events"
	"github.com/anuchandy/coffeemaker/hardwareAPI"
	"github.com/anuchandy/coffeemaker/machineEvent"
)

// ReliefValveController is a controller for coffee-machine's relief-valve.
type ReliefValveController struct {
	Controller
}

// NewReliefValveController returns ReliefValveController, the controller will be registered for the
// events those are used to set the coffee-maker's relief-valve state (OPEN, CLOSED) using api.
func NewReliefValveController(aggregator *events.Aggregator, api hardwareAPI.CommandAPI) *ReliefValveController {
	controller := ReliefValveController{Controller{Api: api}}

	// ReliefValveController is a subscriber for warmer-plate events
	aggregator.SubscribeFunc(controller.handleWarmerPlateEmptyEvent, machineEvent.WamerPlateEmpty)
	aggregator.SubscribeFunc(controller.handleWarmerPlateNotEmptyEvents, machineEvent.WamerPlatePotNotEmpty,
		machineEvent.WamerPlatePotEmpty)
	return &controller
}

// handleWarmerPlateEmptyEvent is the subscribed event handler for warmer-plate empty event.
func (c *ReliefValveController) handleWarmerPlateEmptyEvent(e events.Event) {
	c.Api.SetReliefValveState(hardwareAPI.VALVE_OPEN)
}

// handleWarmerPlateEmptyEvent is the subscribed event handler for warmer-plate not-empty events.
func (c *ReliefValveController) handleWarmerPlateNotEmptyEvents(e events.Event) {
	c.Api.SetReliefValveState(hardwareAPI.VALVE_CLOSED)
}
