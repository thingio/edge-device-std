package models

type State = string

const (
	DeviceStateConnected    State = "connected"
	DeviceStateReconnecting State = "reconnecting"
	DeviceStateDisconnected State = "disconnected"
	DeviceStateException    State = "exception"

	DriverStateRunning State = "running"
)

type DeviceStatus struct {
	Device      *Device `json:"device"`
	State       State   `json:"state"`
	StateDetail string  `json:"state_detail"`
}

type DriverStatus struct {
	Hello                     bool      `json:"hello"`
	Protocol                  *Protocol `json:"protocol"`
	State                     State     `json:"state"`
	StateDetail               string    `json:"state_detail"`
	HealthCheckIntervalSecond int       `json:"health_check_interval_second"`
}
