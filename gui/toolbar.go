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
	treeModel := conf.GetConfig().HostConfigModel
	tb := ToolBar{
		ButtonStyle: ToolBarButtonImageBeforeText,
		Items: []MenuItem{
			Action{
				AssignTo: &addButton,
				Text:     "新增",
				Image:    "res/add.png",
				// Enabled: Bind("isSpecialMode && enabledCB.Checked"),
				OnTriggered: func() {
					icon, _ := walk.NewBitmapFromFile("res/new.png")
					item := &model.HostConfigItem{Name: "test2", Icon: icon}
					treeModel.Append(item)
					common.Info("Add button appended... len: ", len(treeModel.Roots))
				},
			},
			// Separator{},
			Action{
				Text:  "刷新",
				Image: "res/refresh.png",
				// Enabled: Bind("isSpecialMode && enabledCB.Checked"),
				// OnTriggered: mw.specialAction_Triggered,
			},
			// Separator{},
			Action{
				Text:  "修改",
				Image: "res/pencil.png",
				// Enabled: Bind("isSpecialMode && enabledCB.Checked"),
				// OnTriggered: mw.specialAction_Triggered,
			},
			// Separator{},
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
