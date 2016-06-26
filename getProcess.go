package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/tinycedar/lily/common"
)

func main() {
	processes := getAllProcessIds()
	// common.Info("[getAllProcessIds]: %v", processes)
	for _, v := range processes {
		if v > 0 {
			openProcess(v)
		}
	}
}

func getAllProcessIds() []uint32 {
	now := time.Now()
	defer func(now time.Time) {
		common.Info("time elapsed: %v", time.Since(now))
	}(now)
	procEnumProcesses := syscall.NewLazyDLL("Psapi.dll").NewProc("EnumProcesses")
	var processes = make([]uint32, 1024)
	var cbNeeded uint32
	procEnumProcesses.Call(uintptr(unsafe.Pointer(&processes[0])), uintptr(len(processes)), uintptr(unsafe.Pointer(&cbNeeded)))
	if cbNeeded <= 0 {
		common.Error("Calling EnumProcesses returned empty")
		return nil
	}
	return processes[:cbNeeded/4]
}

func openProcess(pid uint32) uint32 {
	procOpenProcess := syscall.NewLazyDLL("Kernel32.dll").NewProc("OpenProcess")
	var dwDesiredAccess uint32 = 0x0400 | 0x0010
	openPid, _, _ := procOpenProcess.Call(uintptr(unsafe.Pointer(&dwDesiredAccess)), 0, uintptr(pid))
	if openPid > 0 {
		procEnumProcessModules := syscall.NewLazyDLL("Psapi.dll").NewProc("EnumProcessModules")
		var cbNeeded uint32
		var modules = make([]unsafe.Pointer, 10)
		if success, _, _ := procEnumProcessModules.Call(uintptr(openPid), uintptr(unsafe.Pointer(&modules[0])), 10, uintptr(unsafe.Pointer(&cbNeeded))); success > 0 {
			// fmt.Println("modules = ", modules)
			// var processName = "<unknown>------------------------"
			procGetModuleBaseName := syscall.NewLazyDLL("Psapi.dll").NewProc("GetModuleBaseNameA")
			var processName = make([]int8, 10240)
			// var processName = ""
			l, _, _ := procGetModuleBaseName.Call(uintptr(openPid), uintptr(unsafe.Pointer(modules[0])), uintptr(unsafe.Pointer(&processName)), uintptr(1024))
			if l > 0 {
				fmt.Println(processName[0:l])
			}
		}
	}
	return 0
}
