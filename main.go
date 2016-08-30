package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/core"
	"github.com/tinycedar/lily/gui"
)

func main() {
	Init()
	// go core.OpenRegistry()
	go core.FireHostsSwitch()
	gui.InitMainWindow()
}

func Init() {
	pidPath := os.TempDir() + "\\lily.pid"
	if !hasStarted(pidPath) {
		if pidFile, err := os.Create(pidPath); err == nil {
			defer pidFile.Close()
			pidFile.WriteString(fmt.Sprint(os.Getpid()))
		} else {
			common.Error("Fail to create pid file, pidPath = %v", pidPath)
		}
	} else {
		common.Info("Already started...")
	}
}

func hasStarted(pidPath string) bool {
	if pidFile, err := os.Open(pidPath); err == nil {
		defer pidFile.Close()
		if bytes, err := ioutil.ReadAll(pidFile); err == nil {
			pid, _ := strconv.Atoi(string(bytes))
			_, err := os.FindProcess(pid)
			return err == nil
		}
	}
	return false
}
