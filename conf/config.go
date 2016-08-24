package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/model"
)

// Config loads and represents config.json
var Config *config

type config struct {
	CurrentHostIndex int
	HostConfigModel  *model.HostConfigModel `json:"-"`
}

func init() {
	common.Init()
	Load()
}

// Load invoked on startup and refresh button clicked
func Load() {
	loadConfig()
	loadHosts()
}

func loadConfig() {
	if startup := Config == nil; startup {
		Config = new(config)
	} else {
		Config.HostConfigModel.RemoveAll()
	}
	Config.HostConfigModel = model.NewHostConfigModel()
	if bytes, err := ioutil.ReadFile("conf/config.json"); err != nil || json.Unmarshal(bytes, Config) != nil {
		common.Error("Fail to read and unmarshal config.json: %v", err)
	}
}

func loadHosts() {
	hostsDir, err := os.Open("conf/hosts")
	if err != nil {
		//TODO define error and notify user
		common.Error("Fail to open conf/hosts directory: %v", err)
		panic(err)
	}
	defer hostsDir.Close()
	hosts, err := hostsDir.Readdir(-1)
	if err != nil {
		//TODO define error and notify user
		common.Error("Fail to read files of conf/hosts directory: %v", err)
		panic(err)
	}
	index := Config.CurrentHostIndex
	for i, f := range hosts {
		temp := strings.Split(f.Name(), ".hosts")
		if len(temp) > 0 {
			icon := common.IconMap[common.ICON_NEW]
			if i == index {
				icon = common.IconMap[common.ICON_OPEN]
			}
			Config.HostConfigModel.Append(&model.HostConfigItem{Name: temp[0], Icon: icon})
		}
	}
	if index < 0 || index >= len(Config.HostConfigModel.Roots) {
		Config.CurrentHostIndex = -1
	}
}
