package gui

import (
	"github.com/lxn/walk"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/core"
)

func newNotify() {
	var err error
	context.notifyIcon, err = walk.NewNotifyIcon()
	if err != nil {
		common.Error("Error invoking NewNotifyIcon: %v", err)
	}
	icon, _ := walk.NewIconFromFile("res/lily.ico")
	if err := context.notifyIcon.SetIcon(icon); err != nil {
		common.Error("Error setting notify icon: %v", err)
	}
	if err := context.notifyIcon.SetToolTip("Click for info or use the context menu to exit."); err != nil {
		common.Error("Fail to set tooltip: %v", err)
	}
	f := func() {
		if !context.mw.Visible() {
			context.mw.Show()
		} else {
			// context.mw.SwitchToThisWindow()
		}
	}
	go core.Triggered(f)
	context.notifyIcon.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			f()
		}
		// if err := context.notifyIcon.ShowCustom(
		// 	"Walk NotifyIcon Example",
		// 	"There are multiple ShowX methods sporting different icons."); err != nil {
		// 	common.Error("Fail to show custom notify: %v", err)
		// }
	})
	exitAction := walk.NewAction()
	if err := exitAction.SetText("退出"); err != nil {
		common.Error("Error setting exitAction text: %v", err)
	}
	exitAction.Triggered().Attach(func() {
		context.notifyIcon.Dispose()
		// os.Exit(-1)
		walk.App().Exit(0)
	})
	if err := context.notifyIcon.ContextMenu().Actions().Add(exitAction); err != nil {
		common.Error("Error Adding exitAction: %v", err)
	}
	if err := context.notifyIcon.SetVisible(true); err != nil {
		common.Error("Error setting notify visible: %v", err)
	}
	// if err := context.notifyIcon.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again."); err != nil {
	// 	common.Error("Error showing info: %v", err)
	// }
}
