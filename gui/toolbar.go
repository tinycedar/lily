package gui

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/conf"
	"github.com/tinycedar/lily/model"
)

func newToolBar() ToolBar {
	tb := ToolBar{
		ButtonStyle: ToolBarButtonImageBeforeText,
		Items: []MenuItem{
			Action{
				AssignTo: &(context.addButton),
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
				AssignTo: &(context.deleteButton),
				Text:     "删除",
				Image:    "res/delete.png",
				// Enabled: Bind("isSpecialMode && enabledCB.Checked"),
				OnTriggered: func() {
					if context.treeView.CurrentItem() == nil {
						walk.MsgBox(context.mw, "删除hosts", "请选择左边列表后再删除", walk.MsgBoxIconInformation)
						context.deleteButton.SetEnabled(false)
						return
					}
					current := context.treeView.CurrentItem().(*model.HostConfigItem)
					message := fmt.Sprintf("确定要删除hosts '%s'?", current.Text())
					ret := walk.MsgBox(context.mw, "删除hosts", message, walk.MsgBoxYesNo)
					if ret == win.IDYES {
						common.Info("Deleting... %v", ret)
						// conf.Config.HostConfigModel.RemoveAll()
					}
					// 	var dlg *walk.Dialog
					// 	Dialog{
					// 		AssignTo: &dlg,
					// 		Title:    "hello",
					// 		MinSize:  Size{300, 300},
					// 		Layout:   VBox{},
					// 		Children: []Widget{
					// 			Composite{
					// 				Layout: HBox{},
					// 				Children: []Widget{
					// 					HSpacer{},
					// 					PushButton{
					// 						// AssignTo: &acceptPB,
					// 						Text: "OK",
					// 						OnClicked: func() {
					// 							// if err := db.Submit(); err != nil {
					// 							// 	log.Print(err)
					// 							// 	return
					// 							// }

					// 							dlg.Accept()
					// 						},
					// 					},
					// 					PushButton{
					// 						// AssignTo:  &cancelPB,
					// 						Text:      "Cancel",
					// 						OnClicked: func() { dlg.Cancel() },
					// 					},
					// 				},
					// 			},
					// 		},
					// 	}.Run(context.mw)
					// 	dlg.Accept()
				},
			},
		},
	}
	return tb
}
