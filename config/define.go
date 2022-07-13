package config

import (
	"flag"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/thingio/edge-device-std/errors"
	"sync"
)

const (
	EnvPrefix  = "eds"
	FilePath   = "etc"
	FileName   = "config"
	FileFormat = "yaml"
)

var cfg = new(Configuration)
var once = sync.Once{}

type Configuration struct {
	DriverOptions  DriverOptions     `json:"driver" yaml:"driver"`
	ManagerOptions ManagerOptions    `json:"manager" yaml:"manager"`
	LogOptions     LogOptions        `json:"log" yaml:"log"`
	MessageBus     MessageBusOptions `json:"msgbus" yaml:"msgbus"`
}

func (c *Configuration) Check() error {
	if err := c.ManagerOptions.Check(); err != nil {
		return err
	}
	return nil
}

func NewConfiguration() (*Configuration, errors.EdgeError) {
	var err error

	once.Do(func() {
		// read the configuration file path
		var configPath string
		flag.StringVar(&configPath, "cp", FilePath, "config file path, e.g. \"/etc\"")
		var configName string
		flag.StringVar(&configName, "cn", FileName, "config file name, e.g. \"config\", excluding the suffix")
		flag.Parse()

		viper.SetEnvPrefix(EnvPrefix)
		viper.AutomaticEnv()
		viper.AddConfigPath(configPath)
		viper.SetConfigName(configName)
		viper.SetConfigType(FileFormat)

		if err = viper.ReadInConfig(); err != nil {
			return
		}

		if err = viper.Unmarshal(cfg, func(dc *mapstructure.DecoderConfig) {
			dc.TagName = FileFormat
		}); err != nil {
			return
		}

		if err = cfg.Check(); err != nil {
			return
		}
	})

	if err != nil {
		return nil, errors.Configuration.Cause(err, "fail to unmarshal the configuration file")
	}
	return cfg, nil
}
