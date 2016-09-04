package main

import (
	"github.com/tinycedar/lily/core"
	"github.com/tinycedar/lily/gui"
)

func main() {
	core.GuaranteeSingleProcess()
	go core.FireHostsSwitch()
	gui.InitMainWindow()
}
