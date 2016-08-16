package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/conf"
	"github.com/tinycedar/lily/model"
)

func newToolBar(treeView *walk.TreeView) ToolBar {
	var addButton *walk.Action
	tb := ToolBar{
		ButtonStyle: ToolBarButtonImageBeforeText,
		Items: []MenuItem{
			Action{
				AssignTo: &addButton,
				Text:     "新增",
				Image:    "res/add.png",
				// Enabled: Bind("isSpecialMode && enabledCB.Checked"),
				OnTriggered: func() {
					item := &model.HostConfigItem{Name: "aaaaaa", Icon: common.IconMap[common.ICON_NEW]}
					conf.Config.HostConfigModel.Append(item)
				},
			},
			//FIXME 去除刷新按钮是因为点击以后， 双击hosts不再生效
			// Action{
			// 	Text:  "刷新",
			// 	Image: "res/refresh.png",
			// 	// Enabled: Bind("isSpecialMode && enabledCB.Checked"),
			// 	OnTriggered: func() {
			// 		conf.Load()
			// 		// item := &model.HostConfigItem{Name: "rrrr", Icon: common.IconMap[common.ICON_NEW]}
			// 		// conf.Config.HostConfigModel.Append(item)
			// 	},
			// },
			// Action{
			// 	Text:  "修改",
			// 	Image: "res/pencil.png",
			// 	// Enabled: Bind("isSpecialMode && enabledCB.Checked"),
			// 	// OnTriggered: mw.specialAction_Triggered,
			// },
			Action{
				Text:  "删除",
				Image: "res/delete.png",
				// Enabled: Bind("isSpecialMode && enabledCB.Checked"),
				// OnTriggered: mw.specialAction_Triggered,
			},
		},
	}
	return tb
}
