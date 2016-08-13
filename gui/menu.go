package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tinycedar/lily/common"
)

func newMenuItems(mw *walk.MainWindow) []MenuItem {
	return []MenuItem{
		Menu{
			Text: "&File",
			Items: []MenuItem{
				Action{
					// AssignTo:    &openAction,
					Text: "&Open",
					// Image:    "../img/open.png",
					// Enabled:  Bind("enabledCB.Checked"),
					// Visible:  Bind("!openHiddenCB.Checked"),
					Shortcut: Shortcut{walk.ModControl, walk.KeyO},
					// OnTriggered: mw.openAction_Triggered,
				},
				Separator{},
				Action{
					Text: "Delete",
					// OnTriggered: func() { mw.Close() },
				},
			},
		},
		Menu{
			Text: "&Help",
			Items: []MenuItem{
				Action{
					// AssignTo:    &showAboutBoxAction,
					Text: "About",
					// OnTriggered: mw.showAboutBoxAction_Triggered,
					OnTriggered: func() {
						// walk.MsgBox(mw, "About", "Developed by @tinycedar", walk.MsgBoxIconInformation)
						common.Info("mw: ", mw == nil)
					},
				},
			},
		},
	}
}
