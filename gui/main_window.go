package gui

import (
	"strings"

	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/model"
)

func InitMainWindow() {
	var mw *walk.MainWindow
	var inTE, outTE *walk.TextEdit
	var treeView *walk.TreeView
	treeModel, err := model.NewHostConfigModel()
	if err != nil {
		common.Error("Error creating host config model: ", err)
		panic(err)
	}
	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Lily",
		MinSize:  Size{800, 500},
		Layout:   VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TreeView{
						AssignTo: &treeView,
						Model:    treeModel,
						OnCurrentItemChanged: func() {
							common.Info("changed...")
							// dir := treeView.CurrentItem().(*Directory)
							// if err := tableModel.SetDirPath(dir.Path()); err != nil {
							// 	walk.MsgBox(
							// 		mainWindow,
							// 		"Error",
							// 		err.Error(),
							// 		walk.MsgBoxOK|walk.MsgBoxIconError)
							// }
						},
					},
					TextEdit{AssignTo: &inTE},
				},
			},
			PushButton{
				Text: "SCREAM",
				OnClicked: func() {
					outTE.SetText(strings.ToUpper(inTE.Text()))
				},
			},
		},
	}).Create(); err != nil {
		common.Error("Error creating main window: ", err)
		os.Exit(-1)
	}
	if icon, err := walk.NewIconFromFile("res/lily.ico"); err != nil {
		common.Error("Error loading icon for main window: ", err)
	} else {
		mw.SetIcon(icon)
	}
	mw.Run()
}
