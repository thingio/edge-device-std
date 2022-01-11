package config

type DriverOptions struct {
	DriverHealthCheckIntervalSecond   int  `json:"driver_health_check_interval_second" yaml:"driver_health_check_interval_second"`
	DeviceHealthCheckIntervalSecond   int  `json:"device_health_check_interval_second" yaml:"device_health_check_interval_second"`
	DeviceAutoReconnect               bool `json:"device_auto_reconnect" yaml:"device_auto_reconnect"` // TODO reconnect automatically by the driver framework
	DeviceAutoReconnectIntervalSecond int  `json:"device_auto_reconnect_interval_second" yaml:"device_auto_reconnect_interval_second"`
	DeviceAutoReconnectMaxRetries     int  `json:"device_auto_reconnect_max_retries" yaml:"device_auto_reconnect_max_retries"`
}

type ManagerOptions struct {
	HTTP struct {
		Port int `json:"port" yaml:"port"`
	} `json:"http" yaml:"http"`
}

type LogOptions struct {
	Path    string `yaml:"path" json:"path"`
	Level   string `yaml:"level" json:"level" default:"info" validate:"regexp=^(info|debug|warn|error)$"`
	Format  string `yaml:"format" json:"format" default:"text" validate:"regexp=^(text|json)$"`
	Console bool   `yaml:"console" json:"console" default:"false"`
	Age     struct {
		Max int `yaml:"max" json:"max" default:"15" validate:"min=1"`
	} `yaml:"age" json:"age"` // days
	Size struct {
		Max int `yaml:"max" json:"max" default:"50" validate:"min=1"`
	} `yaml:"size" json:"size"` // in MB
	Backup struct {
		Max int `yaml:"max" json:"max" default:"15" validate:"min=0"`
	} `yaml:"backup" json:"backup"`
}
