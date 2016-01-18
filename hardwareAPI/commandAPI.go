package hardwareAPI

// CommandAPI is the interface that group API calls to set
// the state of various components of Mark IV Coffee Maker.
type CommandAPI interface {
	/**
	 * This function turns the heating element in the boiler
	 * on or off.
	 */
	SetBoilerState(boilerState BoilerState)

	/**
	 * This function turns the indicator light on or off.
	 * The indicator light should be turned on at the end
	 * of the brewing cycle. It should be turned off when
	 * the user presses the brew button.
	 */
	SetIndicatorState(indicatorState IndicatorState)

	/**
	 * This function opens and closes the pressure-relief
	 * valve. When this valve is closed, steam pressure in
	 * the boiler will force hot water to spray out over
	 * the coffee filter. When the valve is open, the steam
	 * in the boiler escapes into the environment, and the
	 * water in the boiler will not spray out over the filter.
	 */
	SetReliefValveState(reliefValveState ReliefValveState)

	/**
	 * This function turns the heating element in the warmer
	 * plate on or off.
	 */
	SetWarmerPlateState(warmerPlateState WarmerPlateState)
}

// BoilerState is a type represents boiler's heater state.
type BoilerState bool

const (
	BOILER_ON  BoilerState = true
	BOILER_OFF             = false
)

// IndicatorState is a type represents indicator light state.
type IndicatorState bool

const (
	INDICATOR_ON  IndicatorState = true
	INDICATOR_OFF                = false
)

// ReliefValveState is a type represents pressure relief valve state.
type ReliefValveState bool

const (
	VALVE_OPEN   ReliefValveState = true
	VALVE_CLOSED                  = false
)

// WarmerPlateState is a type represents warmer plate state.
type WarmerPlateState bool

const (
	WARMER_ON  WarmerPlateState = true
	WARMER_OFF                  = false
)
