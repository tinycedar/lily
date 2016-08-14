package gui

import (
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
	treeModel := conf.GetConfig().HostConfigModel
	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Lily",
		MinSize:  Size{720, 500},
		Layout:   VBox{},
		// MenuItems: newMenuItems(mw),
		ToolBar: newToolBar(treeView),
		Children: []Widget{
			HSplitter{
				// AlwaysConsumeSpace: true,
				// HandleWidth: 1,
				Children: []Widget{
					TreeView{
						AssignTo: &treeView,
						Model:    treeModel,
						OnCurrentItemChanged: func() {
							// dir := treeView.CurrentItem().(*Directory)
							// if err := tableModel.SetDirPath(dir.Path()); err != nil {
							// 	walk.MsgBox(
							// 		mainWindow,
							// 		"Error",
							// 		err.Error(),
							// 		walk.MsgBoxOK|walk.MsgBoxIconError)
							// }
						},
						StretchFactor: 1,
						OnItemActivated: func() {
							previousIndex := conf.GetConfig().CurrentHostIndex
							for i := 0; i < treeModel.RootCount(); i++ {
								item := treeModel.RootAt(i).(*model.HostConfigItem)
								if item.Text() == treeView.CurrentItem().(*model.HostConfigItem).Text() {
									conf.GetConfig().CurrentHostIndex = i
									icon, _ := walk.NewBitmapFromFile("res/open.png")
									item.Icon = icon
								} else if previousIndex == i {
									icon, _ := walk.NewBitmapFromFile("res/new.png")
									item.Icon = icon
								}
								treeModel.PublishItemChanged(item)
							}
						},
					},
					TextEdit{
						AssignTo:      &hostConfigText,
						StretchFactor: 3,
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
	hostConfigText.SetText("hhhhhhhhhhhhhhhh\r\nwwwwwwwwww\rwwww\nhhhhhhh")
	// toolButton.SetBackground(nil)
	// common.Info("toolbutton: ", toolButton)
	// toolButton
	// common.Info("treeview: ", treeView)
	// bg, err := walk.NewSolidColorBrush(walk.RGB(0, 0, 255))
	// if err != nil {
	// 	common.Info("Error get color: ", err)
	// } else {
	// 	mw.SetBackground(bg)
	// 	common.Info("setting bg", mw.Background())
	// }
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
