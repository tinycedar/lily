package gui

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/conf"
	"github.com/tinycedar/lily/core"
	"github.com/tinycedar/lily/model"
)

func newTreeView(tv **walk.TreeView, hostConfigText **walk.TextEdit) TreeView {
	treeModel := conf.Config.HostConfigModel
	return TreeView{
		AssignTo: tv,
		Model:    treeModel,
		// click
		OnCurrentItemChanged: func() {
			current := (*tv).CurrentItem().(*model.HostConfigItem)
			if bytes, err := ioutil.ReadFile("conf/hosts/" + current.Text() + ".hosts"); err == nil {
				(*hostConfigText).SetText(string(bytes))
			} else {
				(*hostConfigText).SetText("")
			}
		},
		StretchFactor: 1,
		// double click
		OnItemActivated: func() {
			current := (*tv).CurrentItem().(*model.HostConfigItem)
			previousIndex := conf.Config.CurrentHostIndex
			for i := 0; i < treeModel.RootCount(); i++ {
				item := treeModel.RootAt(i).(*model.HostConfigItem)
				if item.Text() == current.Text() {
					conf.Config.CurrentHostIndex = i
					item.Icon = common.IconMap[common.ICON_OPEN]
				} else if previousIndex == i {
					item.Icon = common.IconMap[common.ICON_NEW]
				}
				treeModel.PublishItemChanged(item)
			}
			if err := ioutil.WriteFile("C:/Windows/System32/drivers/etc/hosts", []byte((*hostConfigText).Text()), os.ModeExclusive); err != nil {
				common.Error("Error writing to system hosts file: ", err)
			}
			configJSON, err := json.Marshal(conf.Config)
			if err != nil {
				common.Error("Error marshal json: %v", err)
			} else {
				ioutil.WriteFile("conf/config.json", configJSON, os.ModeExclusive)
			}
			// fire hosts switch
			core.FireHostsSwitch()
		},
	}
}
