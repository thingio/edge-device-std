package config

import (
	"flag"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/thingio/edge-device-std/errors"
)

const (
	EnvPrefix  = "eds"
	FilePath   = "etc"
	FileName   = "config"
	FileFormat = "yaml"
)

type Configuration struct {
	CommonOptions CommonOptions     `json:"common" yaml:"common"`
	MessageBus    MessageBusOptions `json:"msgbus" yaml:"msgbus"`
}

func NewConfiguration() (*Configuration, errors.EdgeError) {
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
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Configuration.Cause(err, "fail to read the configuration file")
	}

	cfg := new(Configuration)
	if err := viper.Unmarshal(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = FileFormat
	}); err != nil {
		return nil, errors.Configuration.Cause(err, "fail to unmarshal the configuration file")
	}

	return cfg, nil
}
