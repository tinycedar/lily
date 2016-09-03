package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"

	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/core"
	"github.com/tinycedar/lily/gui"
)

const pidFilePath = "conf\\lily.pid"

func main() {
	Init()
	go aa()
	go core.FireHostsSwitch()
	gui.InitMainWindow()
}

func aa() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Kill)
	// go func() {
	// 	<-c
	// 	cleanup()
	// 	os.Exit(1)
	// }()
	go func() {
		for {
			sig := <-c
			common.Info("received signal: %v", sig)
		}
	}()
}

func Init() {
	// aa()
	if process := findStartedProcess(); process == nil {
		common.Info("None...")
		if err := ioutil.WriteFile(pidFilePath, []byte(fmt.Sprint(os.Getpid())), os.ModeExclusive); err != nil {
			common.Error("Error writing to system hosts file: ", err)
		}
		// if pidFile, err := os.Create(pidPath); err == nil {
		// 	defer pidFile.Close()
		// 	pidFile.WriteString(fmt.Sprint(os.Getpid()))
		// } else {
		// 	common.Error("Fail to create pid file, pidPath = %v", pidPath)
		// }
	} else {
		common.Info("Found already running process: %v", process.Pid)
		process.Release()
		// process.Kill()
		process.Signal(os.Kill)
		// process.Wait()
		// os.Exit(1)
	}
}

func findStartedProcess() *os.Process {
	if bytes, err := ioutil.ReadFile(pidFilePath); err == nil {
		if pid, err := strconv.Atoi(string(bytes)); err == nil {
			if process, err := os.FindProcess(pid); err == nil {
				if _, ok := core.GetProcessNameMap()[uint32(pid)]; ok {
					return process
				}
			}
		}
	}
	return nil
}
