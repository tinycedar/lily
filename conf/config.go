package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/model"
)

// Config loads and represents config.json
var Config *config

type config struct {
	CurrentHostIndex int `json:"CurrentHostIndex"`
	HostConfigModel  *model.HostConfigModel
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
	startup := Config == nil
	if startup {
		Config = new(config)
		Config.HostConfigModel = model.NewHostConfigModel()
		if bytes, err := ioutil.ReadFile("conf/config.json"); err != nil || json.Unmarshal(bytes, Config) != nil {
			//TODO define error and notify user
			common.Error("Fail to read and unmarshal config.json: ", err)
			panic(err)
		}
	} else {
		Config.HostConfigModel.RemoveAll()
	}
}

func loadHosts() {
	hostsDir, err := os.Open("conf/hosts")
	if err != nil {
		//TODO define error and notify user
		common.Error("Fail to open conf/hosts directory: ", err)
		panic(err)
	}
	defer hostsDir.Close()
	hosts, err := hostsDir.Readdir(-1)
	if err != nil {
		//TODO define error and notify user
		common.Error("Fail to read files of conf/hosts directory: ", err)
		panic(err)
	}
	for i, f := range hosts {
		temp := strings.Split(f.Name(), ".hosts")
		if len(temp) > 0 {
			icon := common.IconMap[common.ICON_NEW]
			if i == Config.CurrentHostIndex {
				icon = common.IconMap[common.ICON_OPEN]
			}
			common.Info("Append hosts " + temp[0])
			fmt.Println("treeModel: ", Config.HostConfigModel)
			Config.HostConfigModel.Append(&model.HostConfigItem{Name: temp[0], Icon: icon})
		}
	}
}
