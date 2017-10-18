package gui

import (
	"io/ioutil"
	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tinycedar/lily/common"
)

func newTextEdit() TextEdit {
	return TextEdit{
		AssignTo:      &(context.hostConfigText),
		StretchFactor: 3,
		Font: Font{
            Family: "Consolas",
            PointSize: 9,
            Bold: false,
        },
		OnKeyUp: func(key walk.Key) {
			current := context.treeView.CurrentItem()
			if current != nil {
				file := "conf/hosts/" + current.Text() + ".hosts"
				if err := ioutil.WriteFile(file, []byte(context.hostConfigText.Text()), os.ModeExclusive); err != nil {
					common.Error("Error writing to system hosts file: ", err)
				}
			}
		},
	}
}
