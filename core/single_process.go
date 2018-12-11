package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tinycedar/lily/common"
)

const pidFilePath = "conf/lily.pid"

var trigger = make(chan bool)

func GuaranteeSingleProcess() {
	if process := findStartedProcess(); process == nil {
		writePidFile(os.Getpid())
		go watchPidFile()
	} else {
		writePidFile(process.Pid)
		os.Exit(1)
	}
}

func findStartedProcess() *os.Process {
	if bytes, err := ioutil.ReadFile(pidFilePath); err == nil {
		if pid, err := strconv.Atoi(string(bytes)); err == nil {
			if process, err := os.FindProcess(pid); err == nil {
				if _, ok := GetProcessNameMap()[uint32(pid)]; ok {
					return process
				}
			}
		}
	}
	return nil
}

func writePidFile(pid int) {
	if err := ioutil.WriteFile(pidFilePath, []byte(fmt.Sprint(pid)), os.ModeExclusive); err != nil {
		common.Error("Error writing to system hosts file: %v", err)
	}
}

func watchPidFile() {
	batcher, err := New(time.Millisecond * 50)
	if err != nil {
		common.Error("Fail to initialize batcher")
		return
	}
	if err := batcher.Add(pidFilePath); err != nil {
		common.Error("Fail to add pid file: %s", pidFilePath)
		return
	}
	for events := range batcher.Events {
		for _, event := range events {
			if event.Op&fsnotify.Write == fsnotify.Write {
				common.Info("modified file: %v", event)
				trigger <- true
				break
			}
		}
	}
}

func Triggered(f func()) {
	for v := range trigger {
		if v {
			f()
		}
	}
}
