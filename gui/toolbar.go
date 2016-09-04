package gui

import (
	"fmt"
	"os"

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
					var dlg *walk.Dialog
					var hostsNameEdit *walk.LineEdit
					Dialog{
						AssignTo: &dlg,
						Title:    "新增",
						MinSize:  Size{300, 150},
						Layout:   VBox{},
						Children: []Widget{
							Composite{
								Layout: Grid{Columns: 2},
								Children: []Widget{
									Label{
										ColumnSpan: 2,
										Text:       "Hosts名字:",
									},
									LineEdit{
										AssignTo:   &hostsNameEdit,
										ColumnSpan: 2,
										// Text:       Bind("PatienceField"),
									},
								},
							},
							HSpacer{},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									PushButton{
										// AssignTo: &acceptPB,
										Text: "确定",
										OnClicked: func() {
											// if err := db.Submit(); err != nil {
											// 	log.Print(err)
											// 	return
											// }
											item := &model.HostConfigItem{Name: hostsNameEdit.Text(), Icon: common.IconMap[common.ICON_NEW]}
											conf.Config.HostConfigModel.Insert(item)
											context.treeView.SetCurrentItem(item)
											dlg.Accept()
										},
									},
									PushButton{
										// AssignTo:  &cancelPB,
										Text:      "取消",
										OnClicked: func() { dlg.Cancel() },
									},
								},
							},
						},
					}.Run(context.mw)
				},
			},
			//FIXME 去除刷新按钮是因为点击以后， 双击hosts不再生效
			// Action{
			// 	Text:  "刷新",
			// 	Image: "res/refresh.png",
			// 	// Enabled: Bind("isSpecialMode && enabledCB.Checked"),
			// 	OnTriggered: func() {
			// 		conf.Load()
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
						if !conf.Config.HostConfigModel.Remove(current) {
							common.Error("Fail to remove current item: %v", current.Text())
							// TODO notify user
							return
						}
						if context.treeView.Model().RootCount() > 0 {
							context.treeView.SetCurrentItem(context.treeView.Model().RootAt(0))
						}
						file := "conf/hosts/" + current.Text() + ".hosts"
						if err := os.Remove(file); err != nil {
							common.Error("Fail to delete file: %s", file)
							// TODO notify user
							return
						}
						common.Info("Succeed to delete file: %s", file)
					}
				},
			},
		},
	}
	return tb
}
