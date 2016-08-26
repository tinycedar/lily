package core

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/tinycedar/lily/common"

	// "github.com/tinycedar/lily/common"
	"strings"
)

var supportedBowsers = []string{
	"chrome.exe", "firefox.exe", "opera.exe",
	"iexplore.exe", "microsoftedge.exe", "microsoftedgecp.exe",
	"sogouexplorer.exe", "qqbrowser.exe", "360se.exe", "360chrome.exe", "liebao.exe", "maxthon.exe", "ucbrowser.exe"}

func getBrowserProcessMap() map[uint32]string {
	pidMap := make(map[uint32]string) // [pid]processName
	for _, v := range getAllProcessIds() {
		if v > 0 {
			processName := strings.ToLower(strings.Trim(openProcess(v), " "))
			if processName == "" {
				continue
			}
			for _, name := range supportedBowsers {
				if strings.HasSuffix(processName, name) {
					pidMap[v] = name
				}
			}
		}
	}
	return pidMap
}

func getAllProcessIds() []uint32 {
	now := time.Now()
	defer func(now time.Time) {
		// common.Info("time elapsed: %v", time.Since(now))
	}(now)
	procEnumProcesses := syscall.NewLazyDLL("Psapi.dll").NewProc("EnumProcesses")
	var processes = make([]uint32, 1024)
	var cbNeeded uint32
	procEnumProcesses.Call(uintptr(unsafe.Pointer(&processes[0])), uintptr(len(processes)), uintptr(unsafe.Pointer(&cbNeeded)))
	if cbNeeded <= 0 {
		// common.Error("Calling EnumProcesses returned empty")
		return nil
	}
	return processes[:cbNeeded/4]
}

func openProcess(pid uint32) string {
	procOpenProcess := syscall.NewLazyDLL("Kernel32.dll").NewProc("OpenProcess")
	var dwDesiredAccess uint32 = 0x0400 | 0x0010
	openPid, _, _ := procOpenProcess.Call(uintptr(unsafe.Pointer(&dwDesiredAccess)), 0, uintptr(pid))
	if openPid <= 0 {
		common.Info("Fail to open process: Pid = %v", pid)
		return ""
	}
	defer closeHandle(openPid)
	return getProcessName(openPid)
}

func closeHandle(openPid uintptr) {
	procCloseHandle := syscall.NewLazyDLL("Kernel32.dll").NewProc("CloseHandle")
	ret, _, _ := procCloseHandle.Call(uintptr(openPid))
	if ret <= 0 {
		common.Info("Fail to close handle: ret = %v", ret)
	}
}

func getProcessName(pid uintptr) string {
	defer metrics("getProcessName")(time.Now())
	procQueryFullProcessImageName := syscall.NewLazyDLL("Kernel32.dll").NewProc("QueryFullProcessImageNameA")
	var cbNeeded uint32 = 1024
	var processName = make([]byte, cbNeeded)
	procQueryFullProcessImageName.Call(uintptr(pid), 0, uintptr(unsafe.Pointer(&processName[0])), uintptr(unsafe.Pointer(&cbNeeded)))
	return string(processName[0:cbNeeded])
}

func getProcessName2(pid uintptr) string {
	defer metrics("getProcessName2")(time.Now())
	procEnumProcessModules := syscall.NewLazyDLL("Psapi.dll").NewProc("EnumProcessModules")
	var cbNeeded uint32
	var modules = make([]unsafe.Pointer, 10)
	if success, _, _ := procEnumProcessModules.Call(uintptr(pid), uintptr(unsafe.Pointer(&modules[0])), 10, uintptr(unsafe.Pointer(&cbNeeded))); success > 0 {
		procGetModuleBaseName := syscall.NewLazyDLL("Psapi.dll").NewProc("GetModuleBaseNameA")
		var processName = make([]byte, 1024)
		cbNeeded = 20
		l, _, _ := procGetModuleBaseName.Call(uintptr(pid), uintptr(unsafe.Pointer(modules[0])), uintptr(unsafe.Pointer(&processName[0])), uintptr(cbNeeded))
		if l > 0 {
			return string(processName[0:cbNeeded])
		}
	}
	return ""
}

func metrics(funcName string) func(now time.Time) {
	return func(now time.Time) {
		// common.Info("Processing [%v] costs %v\n", funcName, time.Since(now))
	}
}
