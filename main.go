package main

import (
	"github.com/tinycedar/lily/core"
	"github.com/tinycedar/lily/gui"
)

func main() {
	go core.OpenRegistry()
	go core.NewSystemHostsWatcher()
	gui.InitMainWindow()
}
