package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	i := 0
	for _, v := range getAllProcessesIds() {
		i++
		fmt.Println(v)
	}
	fmt.Println("i= ", i)
}

func getAllProcessesIds() []uint32 {
	now := time.Now()
	defer fmt.Println("time elapsed: ", time.Since(now))
	return enumProcesses(1024)
}

func enumProcesses(val uint32) []uint32 {
	procEnumProcesses := syscall.NewLazyDLL("Psapi.dll").NewProc("EnumProcesses")
	var aProcesses = make([]uint32, val)
	var cbNeeded uint32
	procEnumProcesses.Call(uintptr(unsafe.Pointer(&aProcesses[0])), uintptr(val), uintptr(unsafe.Pointer(&cbNeeded)))
	fmt.Println("cbNeeded = ", cbNeeded)
	if cbNeeded > 0 {
		return aProcesses[:cbNeeded/4]
	}
	return enumProcesses(val * 2)
}
