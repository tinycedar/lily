package gui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/conf"
	"github.com/tinycedar/lily/model"
)

func InitMainWindow() {
	var mw *walk.MainWindow
	var hostConfigText *walk.TextEdit
	var treeView *walk.TreeView
	treeModel := conf.Config.HostConfigModel
	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Lily",
		MinSize:  Size{720, 500},
		Layout:   VBox{},
		// MenuItems: newMenuItems(mw),
		ToolBar: newToolBar(treeView),
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TreeView{
						AssignTo: &treeView,
						Model:    treeModel,
						OnCurrentItemChanged: func() {
							item := treeView.CurrentItem().(*model.HostConfigItem)
							bytes, err := ioutil.ReadFile("conf/hosts/" + item.Text() + ".hosts")
							if err == nil {
								hostConfigText.SetText(fmt.Sprintf("%s", bytes))
							} else {
								hostConfigText.SetText("")
							}
							// 	walk.MsgBox(
							// 		mainWindow,
							// 		"Error",
							// 		err.Error(),
							// 		walk.MsgBoxOK|walk.MsgBoxIconError)
							// }
						},
						StretchFactor: 1,
						OnItemActivated: func() {
							current := treeView.CurrentItem().(*model.HostConfigItem)
							previousIndex := conf.Config.CurrentHostIndex
							for i := 0; i < treeModel.RootCount(); i++ {
								item := treeModel.RootAt(i).(*model.HostConfigItem)
								if item.Text() == current.Text() {
									conf.Config.CurrentHostIndex = i
									icon, _ := walk.NewBitmapFromFile("res/open.png")
									item.Icon = icon
								} else if previousIndex == i {
									icon, _ := walk.NewBitmapFromFile("res/new.png")
									item.Icon = icon
								}
								treeModel.PublishItemChanged(item)
							}
							if err := ioutil.WriteFile("C:/Windows/System32/drivers/etc/hosts", []byte(hostConfigText.Text()), os.ModeExclusive); err != nil {
								common.Error("Error writing to system hosts file: ", err)
							}
							configJSON, err := json.Marshal(conf.Config)
							if err != nil {
								common.Error("Error marshal json: ", err)
							} else {
								common.Info("configJSON: ", string(configJSON))
								ioutil.WriteFile("conf/config.json", configJSON, os.ModeExclusive)
							}
						},
					},
					TextEdit{
						AssignTo:      &hostConfigText,
						StretchFactor: 3,
						OnKeyUp: func(key walk.Key) {
							common.Info("============================ Key up =================================")
						},
					},
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
	setXY(mw)

	item := treeModel.RootAt(conf.Config.CurrentHostIndex)
	if bytes, err := ioutil.ReadFile("conf/hosts/" + item.Text() + ".hosts"); err != nil {
		common.Error("Error reading host config: ", err)
	} else {
		hostConfigText.SetText(fmt.Sprintf("%s", bytes))
	}
	bg, err := walk.NewSolidColorBrush(walk.RGB(218, 223, 230))
	if err != nil {
		common.Error("Error new color brush", err)
	} else {
		treeView.SetBackground(bg)
	}
	mw.Run()
}

func setXY(mw *walk.MainWindow) {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)
	mw.SetX((int(win.GetDeviceCaps(hDC, win.HORZRES)) - mw.Width()) / 2)
	mw.SetY((int(win.GetDeviceCaps(hDC, win.VERTRES)) - mw.Height()) / 2)
}
