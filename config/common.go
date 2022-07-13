package config

import (
	"fmt"
	"time"
)

type (
	MetaStoreType = string
	DataStoreType = string
)

const (
	MetaStoreTypeFile MetaStoreType = "file"

	DataStoreTypeInfluxDB DataStoreType = "influxdb"
	DataStoreTypeTDengine DataStoreType = "tdengine"

	MinBatchSize = 1
	MaxBatchSize = 10000
)

type DriverOptions struct {
	DriverHealthCheckIntervalSecond   int  `json:"driver_health_check_interval_second" yaml:"driver_health_check_interval_second"`
	DeviceHealthCheckIntervalSecond   int  `json:"device_health_check_interval_second" yaml:"device_health_check_interval_second"`
	DeviceAutoReconnect               bool `json:"device_auto_reconnect" yaml:"device_auto_reconnect"` // TODO reconnect automatically by the driver framework
	DeviceAutoReconnectIntervalSecond int  `json:"device_auto_reconnect_interval_second" yaml:"device_auto_reconnect_interval_second"`
	// The number of retries for automatic reconnection of the device. If it is 0, there is no limit.
	DeviceAutoReconnectMaxRetries int `json:"device_auto_reconnect_max_retries" yaml:"device_auto_reconnect_max_retries"`
}

type ManagerOptions struct {
	HTTP struct {
		Port int `json:"port" yaml:"port"`
	} `json:"http" yaml:"http"`

	MetaStoreOptions *MetaStoreOptions `json:"meta_store" yaml:"meta_store"`
	DataStoreOptions *DataStoreOptions `json:"data_store" yaml:"data_store"`
}

func (o *ManagerOptions) Check() error {
	if o.DataStoreOptions == nil {
		return nil
	}
	if err := o.DataStoreOptions.Check(); err != nil {
		return err
	}
	return nil
}

type MetaStoreOptions struct {
	Type MetaStoreType `json:"type" yaml:"type"`
	File *FileOptions  `json:"file" yaml:"file"`
}

type FileOptions struct {
	Path string `json:"path" yaml:"path"`
}

type DataStoreOptions struct {
	Type DataStoreType `json:"type" yaml:"type"`

	// Common for All DBs
	// https://docs.taosdata.com/taos-sql/database#创建数据库
	Database    string        `json:"database" yaml:"database"`
	BatchSize   int           `json:"batch_size" yaml:"batch_size"`

	InfluxDB *InfluxDBOptions `json:"influxdb" yaml:"influxdb"`
	TDengine *TDengineOptions `json:"tdengine" yaml:"tdengine"`
}

func (o *DataStoreOptions) Check() error {
	if o.Database == "" {
		return fmt.Errorf("influxdb database must be specified")
	}

	if o.BatchSize > MaxBatchSize {
		return fmt.Errorf("max-batch-size cannot be large than %d", MaxBatchSize)
	}
	if o.BatchSize < MinBatchSize {
		return fmt.Errorf("batch-size cannot be less than %d", MinBatchSize)
	}

	switch o.Type {
	case DataStoreTypeInfluxDB:
		if o.InfluxDB == nil {
			return fmt.Errorf("the configuration for InfluxDB must be required")
		}
		return o.InfluxDB.Check()
	case DataStoreTypeTDengine:
		if o.TDengine == nil {
			return fmt.Errorf("the configuration for TDengine must be required")
		}
		return o.TDengine.Check()
	default:
		return fmt.Errorf("unsupported datastore type: %s", o.Type)
	}
}

type InfluxDBOptions struct {
	URL              string        `json:"url" yaml:"url"`
	Username         string        `json:"username" yaml:"username"`
	Password         string        `json:"password" yaml:"password"`
	UserAgent        string        `json:"user_agent" yaml:"user_agent"`
	Precision        string        `json:"precision" yaml:"precision"`
	RetentionPolicy  string        `json:"retention_policy" yaml:"retention_policy"`
	WriteConsistency string        `json:"write_consistency" yaml:"write_consistency"`
	Timeout          time.Duration `json:"timeout" yaml:"timeout"`
}

func (o *InfluxDBOptions) Check() error {
	return nil
}

type TDengineOptions struct {
	Schema    string `json:"schema" yaml:"schema"`
	Host      string `json:"host" yaml:"host"`
	Port      uint16 `json:"port" yaml:"port"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	Keep      uint   `json:"keep" yaml:"keep"`
	Days      uint   `json:"days" yaml:"days"`
	Blocks    uint   `json:"blocks" yaml:"blocks"`
	Update    uint8  `json:"update" yaml:"update"`
	Precision string `json:"precision" yaml:"precision"`

	MaxBatchSize int `json:"max_batch_size" yaml:"max_batch_size"`
}

func (o *TDengineOptions) Check() error {
	return nil
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
