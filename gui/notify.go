package gui

import (
	"log"

	"github.com/lxn/walk"
	"github.com/tinycedar/lily/common"
)

func newNotify(mw *walk.MainWindow) {
	ni, err := walk.NewNotifyIcon()
	if err != nil {
		common.Error("Error invoking NewNotifyIcon", err)
	}
	icon, _ := walk.NewIconFromFile("res/lily.ico")
	if err := ni.SetIcon(icon); err != nil {
		common.Error("Error setting notify icon", err)
	}
	if err := ni.SetToolTip("Click for info or use the context menu to exit."); err != nil {
		log.Fatal(err)
	}
	ni.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			if !mw.Visible() {
				mw.Show()
			} else {
				//TODO
			}
		}
		// if err := ni.ShowCustom(
		// 	"Walk NotifyIcon Example",
		// 	"There are multiple ShowX methods sporting different icons."); err != nil {
		// 	log.Fatal(err)
		// }
	})
	exitAction := walk.NewAction()
	if err := exitAction.SetText("退出"); err != nil {
		common.Error("Error setting exitAction text", err)
	}
	exitAction.Triggered().Attach(func() {
		defer ni.Dispose()
		walk.App().Exit(0)
	})
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		common.Error("Error Adding exitAction", err)
	}
	if err := ni.SetVisible(true); err != nil {
		common.Error("Error setting notify visible", err)
	}
	// if err := ni.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again."); err != nil {
	// 	common.Error("Error showing info", err)
	// }
}
