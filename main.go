package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"unicode/utf16"
	"unsafe"

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
				if _, ok := getProcessNameMap()[uint32(pid)]; ok {
					return process
				}
			}
		}
	}
	return nil
}

func getProcessNameMap() map[uint32]string {
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		common.Error("Fail to syscall CreateToolhelp32Snapshot: %v", err)
		return nil
	}
	defer syscall.CloseHandle(snapshot)
	var procEntry syscall.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = syscall.Process32First(snapshot, &procEntry); err != nil {
		common.Error("Fail to syscall Process32First: %v", err)
		return nil
	}
	processNameMap := make(map[uint32]string)
	for {
		processNameMap[procEntry.ProcessID] = parseProcessName(procEntry.ExeFile)
		if err = syscall.Process32Next(snapshot, &procEntry); err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				return processNameMap
			}
			common.Error("Fail to syscall Process32Next: %v", err)
			return nil
		}
	}
}

func parseProcessName(exeFile [syscall.MAX_PATH]uint16) string {
	for i, v := range exeFile {
		if v <= 0 {
			return string(utf16.Decode(exeFile[:i]))
		}
	}
	return ""
}
