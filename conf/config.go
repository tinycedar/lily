package conf

import (
	"github.com/tinycedar/lily/model"
)

var config *Config

func init() {
	// loading from json file
	config = &Config{}
	config.CurrentHostIndex = -1
	config.HostConfigModel = model.NewHostConfigModel()
	//TODO write back
}

type Config struct {
	CurrentHostIndex int
	HostConfigModel  *model.HostConfigModel
}

func GetConfig() *Config {
	return config
}
