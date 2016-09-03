package core

import (
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/tinycedar/lily/common"
)

func GetProcessNameMap() map[uint32]string {
	defer metrics("GetProcessNameMap")(time.Now())
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

func metrics(funcName string) func(now time.Time) {
	return func(now time.Time) {
		common.Info("Processing [%v] costs %v", funcName, time.Since(now))
	}
}
