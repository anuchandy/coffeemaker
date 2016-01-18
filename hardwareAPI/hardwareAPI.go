package hardwareAPI

// HardwareAPI is the interface that groups the API calls to query status
// and set state of various components of the Mark IV Coffee Maker.
type HardwareAPI interface {
	QueryAPI
	CommandAPI
}
