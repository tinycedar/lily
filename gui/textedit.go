package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tinycedar/lily/common"
)

func newTextEdit(te **walk.TextEdit) TextEdit {
	return TextEdit{
		AssignTo:      te,
		StretchFactor: 3,
		OnKeyUp: func(key walk.Key) {
			common.Info("=========== Key up ===========")
		},
	}
}
