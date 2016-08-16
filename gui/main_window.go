package gui

import (
	"io/ioutil"
	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/conf"
)

func InitMainWindow() {
	var mw *walk.MainWindow
	var treeView = new(walk.TreeView)
	var hostConfigText *walk.TextEdit
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
					newTreeView(&treeView, &hostConfigText),
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
	setXY(mw)
	setBackground(treeView)
	icon, _ := walk.NewIconFromFile("res/lily.ico")
	mw.SetIcon(icon)

	currentItem := conf.Config.HostConfigModel.RootAt(conf.Config.CurrentHostIndex)
	if currentItem == nil {
		common.Error("Invalid CurrentHostIndex in config.json, cannot find the specific hosts")
	} else {
		if bytes, err := ioutil.ReadFile("conf/hosts/" + currentItem.Text() + ".hosts"); err != nil {
			common.Error("Error reading host config: ", err)
		} else {
			hostConfigText.SetText(string(bytes))
		}
	}
	mw.Run()
}

func setXY(mw *walk.MainWindow) {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)
	mw.SetX((int(win.GetDeviceCaps(hDC, win.HORZRES)) - mw.Width()) / 2)
	mw.SetY((int(win.GetDeviceCaps(hDC, win.VERTRES)) - mw.Height()) / 2)
}

func setBackground(treeView *walk.TreeView) {
	bg, err := walk.NewSolidColorBrush(walk.RGB(218, 223, 230))
	if err != nil {
		common.Error("Error new color brush", err)
	} else {
		treeView.SetBackground(bg)
	}
}
