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

var context = new(widgetContext)

type widgetContext struct {
	mw             *walk.MainWindow
	treeView       *walk.TreeView
	hostConfigText *walk.TextEdit
	addButton      *walk.Action
	deleteButton   *walk.Action
}

// InitMainWindow initialize main window
func InitMainWindow() {
	if err := (MainWindow{
		AssignTo: &(context.mw),
		Title:    "Lily",
		MinSize:  Size{720, 500},
		Layout:   VBox{},
		// MenuItems: newMenuItems(mw),
		ToolBar: newToolBar(),
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					newTreeView(),
					newTextEdit(),
				},
			},
		},
	}).Create(); err != nil {
		common.Error("Error creating main window: %v", err)
		os.Exit(-1)
	}
	setWindowIcon(context.mw)
	setXY(context.mw)
	newNotify(context.mw)
	setTreeViewBackground(context.treeView)
	showCurrentItem(context.hostConfigText)
	context.mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		*canceled = true
		context.mw.Hide()
	})
	(*context.deleteButton).SetEnabled(false)
	context.mw.Run()
}

func setWindowIcon(mw *walk.MainWindow) {
	icon, _ := walk.NewIconFromFile("res/lily.ico")
	mw.SetIcon(icon)
}

func setXY(mw *walk.MainWindow) {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)
	mw.SetX((int(win.GetDeviceCaps(hDC, win.HORZRES)) - mw.Width()) / 2)
	mw.SetY((int(win.GetDeviceCaps(hDC, win.VERTRES)) - mw.Height()) / 2)
}

func setTreeViewBackground(treeView *walk.TreeView) {
	bg, err := walk.NewSolidColorBrush(walk.RGB(218, 223, 230))
	if err != nil {
		common.Error("Error new color brush: %v", err)
	} else {
		treeView.SetBackground(bg)
	}
}

func showCurrentItem(hostConfigText *walk.TextEdit) {
	model := conf.Config.HostConfigModel
	index := conf.Config.CurrentHostIndex
	if index < 0 || len(model.Roots) == 0 {
		return
	}
	current := model.RootAt(index)
	if current == nil {
		common.Error("Invalid CurrentHostIndex in config.json, cannot find the specific hosts")
	} else {
		if bytes, err := ioutil.ReadFile("conf/hosts/" + current.Text() + ".hosts"); err != nil {
			common.Error("Error reading host config: %v", err)
		} else {
			hostConfigText.SetText(string(bytes))
		}
	}
}
