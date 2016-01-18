package hardwareAPI

// QueryAPI is the interface that group API calls to query the
// status of various components of Mark IV Coffee Maker.
type QueryAPI interface {
	/**
	 * This function returns the status of the boiler switch.
	 * The boiler switch is a float switch that detects if there is
	 * more than 1/2 cup of water in the boiler.
	 */
	GetBoilerStatus() BoilerStatus

	/**
	 * This function returns the status of the brew button. The brew
	 * button is a momentary switch that remembers it's state.
	 * Each call to this function returns the remembered state and
	 * then resets that state to BREW_BUTTON_NOT_PUSHED.
	 *
	 * Thus, even if this function is polled at a very slow rate, it
	 * will still detect when the brew button is pushed.
	 */
	GetBrewButtonStatus() BrewButtonStatus

	/**
	 * This function returns the status of the warmer-plate sensor.
	 * This sensor detects the presence of the pot and whether it
	 * has coffee in it.
	 */
	GetWarmerPlateStatus() WarmerStatus
}

// BoilerStatus is a type represents boiler water status.
type BoilerStatus bool

const (
	BOILER_EMPTY     BoilerStatus = true
	BOILER_NOT_EMPTY              = false
)

// BrewButtonStatus is a type represents brew button status.
type BrewButtonStatus bool

const (
	BREW_BUTTON_PUSHED     BrewButtonStatus = true
	BREW_BUTTON_NOT_PUSHED                  = false
)

// WarmerStatus is a type represents warmer plate's heater status.
type WarmerStatus byte

const (
	WARMER_EMPTY WarmerStatus = iota
	POT_EMPTY
	POT_NOT_EMPTY
)
