package hardwareController

import (
	"github.com/anuchandy/coffeemaker/events"
	"github.com/anuchandy/coffeemaker/hardwareAPI"
	"github.com/anuchandy/coffeemaker/machineEvent"
)

// BoilerController is a controller for coffee-machine's indicator light.
type IndicatorController struct {
	Controller
	boilerEmpty       bool
	brewingInProgress bool
}

// NewIndicatorController returns IndicatorController, the controller will be registered for the
// events those are used to set the coffee-maker's indicator light state (ON, OFF) using api.
func NewIndicatorController(aggregator *events.Aggregator, api hardwareAPI.CommandAPI) *IndicatorController {
	controller := IndicatorController{
		Controller:        Controller{Api: api},
		boilerEmpty:       false,
		brewingInProgress: false,
	}

	// IndicatorController is a subscriber for boiler events
	aggregator.SubscribeFunc(controller.handleBoilerEmptyEvent, machineEvent.BoilerEmpty)
	aggregator.SubscribeFunc(controller.handleBoilerNotEmptyEvent, machineEvent.BoilerNotEmpty)
	// IndicatorController is a subscriber for brew button event BrewButtonPushed
	aggregator.SubscribeFunc(controller.handleBrewButtonPushEvent, machineEvent.BrewButtonPushed)

	return &controller
}

// handleBoilerEmptyEvent is the subscribed event handler for boiler empty event.
func (c *IndicatorController) handleBoilerEmptyEvent(e events.Event) {
	c.boilerEmpty = true
	if c.brewingInProgress {
		c.brewingInProgress = false
		c.Api.SetIndicatorState(hardwareAPI.INDICATOR_ON)
	}
}

// handleBoilerNotEmptyEvent is the subscribed event handler for boiler not-empty event.
func (c *IndicatorController) handleBoilerNotEmptyEvent(e events.Event) {
	c.boilerEmpty = false
}

// handleBrewButtonPushEvent is the subscribed event handler for brew button pushed event.
func (c *IndicatorController) handleBrewButtonPushEvent(e events.Event) {
	if !c.boilerEmpty {
		c.brewingInProgress = true
		c.Api.SetIndicatorState(hardwareAPI.INDICATOR_OFF)
	}
}
