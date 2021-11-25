package config

type CommonOptions struct {
	DriverHealthCheckIntervalSecond   int  `json:"driver_health_check_interval_second" yaml:"driver_health_check_interval_second"`
	DeviceHealthCheckIntervalSecond   int  `json:"device_health_check_interval_second" yaml:"device_health_check_interval_second"`
	DeviceAutoReconnect               bool `json:"device_auto_reconnect" yaml:"device_auto_reconnect"`
	DeviceAutoReconnectIntervalSecond int  `json:"device_auto_reconnect_interval_second" yaml:"device_auto_reconnect_interval_second"`
}
