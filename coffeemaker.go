package coffeemaker

import (
	"github.com/anuchandy/coffeemaker/events"
	"github.com/anuchandy/coffeemaker/hardwareAPI"
	"github.com/anuchandy/coffeemaker/hardwareController"
	"github.com/anuchandy/coffeemaker/machineEvent"
	"time"
)

// The event aggregator for publishing and subscribing for coffee-maker's events.
var agg *events.Aggregator = events.NewAggregator()
// The channel to send stop polling signal to hardware polling go routine.
var abortPoll chan struct{} = make(chan struct{})

// SwitchOn switch-on the coffee-maker.
func SwitchOn(api hardwareAPI.HardwareAPI) {
	// Create the hardware controllers those will subscribe for machine events
	// from event-aggregator.
	createControllers(api)

	// Start the event aggregator [async].
	agg.Start()

	// poll the hardware and publish machine events to event-aggregator [async].
	pollHardware(api)
}

// SwitchOff switch-off the coffee-maker.
func SwitchOff() {
	abortPoll <- struct{}{}
	agg.Stop()
}

// createControllers creates controllers for various components of the coffee maker
// those are subscribed for events required to control the components.
func createControllers(api hardwareAPI.CommandAPI) {
	hardwareController.NewBoilerController(agg, api)
	hardwareController.NewIndicatorController(agg, api)
	hardwareController.NewReliefValveController(agg, api)
	hardwareController.NewWarmerPlateController(agg, api)
}

// pollHardware polls the coffee maker's hardware and publishes status of various
// components as events.
//
// pollHardware stops the polling if it receives signal in abortPoll channel.
func pollHardware(api hardwareAPI.QueryAPI) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				publishEvents(api)
			case <-abortPoll:
				ticker.Stop()
				return
			}
		}
	}()
}

// publishEvents publishes the current state of coffee-maker.
func publishEvents(api hardwareAPI.QueryAPI) {
	agg.Publish(toBoilerEvent(api.GetBoilerStatus()))
	agg.Publish(toBrewButtonEvent(api.GetBrewButtonStatus()))
	agg.Publish(toWarmerPlateEvent(api.GetWarmerPlateStatus()))
}

// toBoilerEvent returns the MachineEvent corresponds to the boiler status s.
func toBoilerEvent(s hardwareAPI.BoilerStatus) *machineEvent.MachineEvent {
	if s == hardwareAPI.BOILER_EMPTY {
		return machineEvent.BoilerEmpty
	}
	return machineEvent.BoilerNotEmpty
}

// toBrewButtonEvent returns the MachineEvent corresponds to the brew button status s.
func toBrewButtonEvent(s hardwareAPI.BrewButtonStatus) *machineEvent.MachineEvent {
	if s == hardwareAPI.BREW_BUTTON_PUSHED {
		return machineEvent.BrewButtonPushed
	}
	return machineEvent.BrewButtonNotPushed
}

// toWarmerPlateEvent returns the MachineEvent corresponds to the warmer-plate status s.
func toWarmerPlateEvent(s hardwareAPI.WarmerStatus) *machineEvent.MachineEvent {
	if s == hardwareAPI.WARMER_EMPTY {
		return machineEvent.WamerPlateEmpty
	}

	if s == hardwareAPI.POT_EMPTY {
		return machineEvent.WamerPlatePotEmpty
	}
	return machineEvent.WamerPlatePotNotEmpty
}
