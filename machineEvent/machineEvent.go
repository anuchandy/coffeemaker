package machineEvent

import (
	"strconv"
)

// MachineEvent represents coffee-maker hardware event.
type MachineEvent int64

// NewMachineEvent returns MachineEvent with the given id.
func NewMachineEvent(id int64) *MachineEvent {
	e := MachineEvent(id)
	return &e
}

// There is more than 1/2 cup of water in the boiler.
var BoilerEmpty = NewMachineEvent(0)
// There is no water in the boiler or water level is less than 1/2 cup.
var BoilerNotEmpty = NewMachineEvent(1)
// There is no pot in the warmer plate.
var WamerPlateEmpty = NewMachineEvent(2)
// Warmer plate has pot and it is empty.
var WamerPlatePotEmpty = NewMachineEvent(3)
// Warmer plate has pot and it is filled with coffee.
var WamerPlatePotNotEmpty = NewMachineEvent(4)
// Brew button is not pushed.
var BrewButtonNotPushed = NewMachineEvent(5)
// Brew button is pushed.
var BrewButtonPushed = NewMachineEvent(6)

var machineEvents = [...]string{
	0: "BoilerEmpty",
	1: "BoilerNotEmpty",
	2: "WamerPlateEmpty",
	3: "WamerPlatePotEmpty",
	4: "WamerPlatePotNotEmpty",
	5: "BrewButtonNotPushed",
	6: "BrewButtonPushed",
}

// Event method to satisfies event interface.
func (m *MachineEvent) Event() string {
	if 0 <= *m && int(*m) < len(machineEvents) {
		return machineEvents[*m]
	}

	return strconv.FormatInt(int64(*m), 10)
}

// String method to satisfies Stringer interface.
func (m *MachineEvent) String() string {
	return m.Event()
}
