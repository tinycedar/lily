package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tinycedar/lily/common"
)

func newTextEdit() TextEdit {
	return TextEdit{
		AssignTo:      &(context.hostConfigText),
		StretchFactor: 3,
		OnKeyUp: func(key walk.Key) {
			common.Info("=========== Key up ===========")
		},
	}
}
