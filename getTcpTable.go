package main

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

var bigEndian = binary.LittleEndian

type ClassReader struct {
	bytecode []byte
}

func NewClassReader(bytecode []byte) *ClassReader {
	return &ClassReader{bytecode: bytecode}
}

func (this *ClassReader) ReadUint32() uint32 {
	value := bigEndian.Uint32(this.bytecode[:4])
	this.bytecode = this.bytecode[4:]
	return value
}

func (this *ClassReader) ReadBytes(len int) []byte {
	bytes := this.bytecode[:len]
	this.bytecode = this.bytecode[len:]
	return bytes
}

func (this *ClassReader) ReadIp(bytes []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", bytes[0], bytes[1], bytes[2], bytes[3])
}

func (this *ClassReader) ReadPort(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes[0:2])
}

type (
	DWORD uint32
)

type TCP_CONNECTION_OFFLOAD_STATE uint32

const (
	TcpConnectionOffloadStateInHost     TCP_CONNECTION_OFFLOAD_STATE = 0
	TcpConnectionOffloadStateOffloading TCP_CONNECTION_OFFLOAD_STATE = 1
	TcpConnectionOffloadStateOffloaded  TCP_CONNECTION_OFFLOAD_STATE = 2
	TcpConnectionOffloadStateUploading  TCP_CONNECTION_OFFLOAD_STATE = 3
	TcpConnectionOffloadStateMax        TCP_CONNECTION_OFFLOAD_STATE = 4
)

type MIB_TCPROW2 struct {
	dwState        DWORD
	dwLocalAddr    DWORD
	dwLocalPort    DWORD
	dwRemoteAddr   DWORD
	dwRemotePort   DWORD
	dwOwningPid    DWORD
	dwOffloadState TCP_CONNECTION_OFFLOAD_STATE
}

func (r *MIB_TCPROW2) displayIP(val DWORD) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(val), byte(val>>8), byte(val>>16), val>>24)
}

func (r *MIB_TCPROW2) displayPort(val DWORD) uint16 {
	int32Val := uint32(val)
	return binary.BigEndian.Uint16([]byte{byte(int32Val), byte(int32Val >> 8)})
}

func newTCPRow(r *ClassReader) *MIB_TCPROW2 {
	return &MIB_TCPROW2{DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), TCP_CONNECTION_OFFLOAD_STATE(r.ReadUint32())}
}

type MIB_TCPTABLE2 struct {
	dwNumEntries DWORD
	table        []*MIB_TCPROW2
}

func (t *MIB_TCPTABLE2) String() string {
	fmt.Println("================  tcp table ======================= ", t.dwNumEntries)
	for i := uint32(0); i < uint32(t.dwNumEntries); i++ {
		row := t.table[i]
		fmt.Println(row, "\t", row.displayIP(row.dwRemoteAddr), ":", row.displayPort(row.dwRemotePort))
	}
	fmt.Println("================  tcp table end =======================")
	return "======================================="
}

func newTCPTable(r *ClassReader) *MIB_TCPTABLE2 {
	t := &MIB_TCPTABLE2{}
	t.dwNumEntries = DWORD(r.ReadUint32())
	table := make([]*MIB_TCPROW2, t.dwNumEntries)
	for i := uint32(0); i < uint32(t.dwNumEntries); i++ {
		table[i] = newTCPRow(r)
	}
	t.table = table
	return t
}

func main() {
	call := syscall.NewLazyDLL("Iphlpapi.dll")
	getTcpTable2 := call.NewProc("GetTcpTable2")
	var n uint32 = 0
	table := &MIB_TCPTABLE2{}
	r1, _, _ := getTcpTable2.Call(uintptr(unsafe.Pointer(table)), uintptr(unsafe.Pointer(&n)), 1)
	// if r1 != syscall.ERROR_INSUFFICIENT_BUFFER {
	// something bad happened; use syscall.Error(r1) to diagnose
	fmt.Println(r1)
	// }
	b := make([]byte, n)
	r2, _, _ := getTcpTable2.Call(uintptr(unsafe.Pointer(&b[0])), uintptr(unsafe.Pointer(&n)), 1)
	if r2 != 0 {
		// something bad happened; use syscall.Error(r1) to diagnose
		fmt.Println(r2)
	}
	const (
		BING string = "202.89.233.103" // netstat -ano | findstr 202.89.233.103
		// kaola
	)
	testRemoteAddr := BING
	table = newTCPTable(NewClassReader(b))
	// fmt.Println(table)
	for i := uint32(0); i < uint32(table.dwNumEntries); i++ {
		row := table.table[i]
		remoteAddr := row.displayIP(row.dwRemoteAddr)
		if remoteAddr == "0.0.0.0" || remoteAddr == "127.0.0.1" || remoteAddr != testRemoteAddr {
			continue
		}
		if row.dwOwningPid > 0 && remoteAddr == testRemoteAddr {
			fmt.Println("row: ", row)
			setTcpEntry := call.NewProc("SetTcpEntry")
			row.dwState = 12
			r, _, _ := setTcpEntry.Call(uintptr(unsafe.Pointer(row)))
			fmt.Println("setTcpEntry: ", r)
		}
		// fmt.Printf("\t%-6d\t%s:%-16d\t%s:%-16d\t%d\t%d\n", state, r.ReadIp(localAddrBytes), r.ReadPort(localPortBytes), remoteAddr, remotePort, pid, offLoadState)
	}
	// loopCount := reader.ReadUint32()
	// fmt.Println("count = ", loopCount)
	// fmt.Printf("\t%s\t%16s\t%16s\t\t\t%s\t%s\n", "State", "Local Address", "Foreign Address", "PID", "Off Load State")
	// for i := uint32(0); i < loopCount; i++ {
	// 	bb := reader.ReadBytes(28)
	// 	// fmt.Println("bytes: ", bb)
	// 	r := NewClassReader(bb)
	// 	state := r.ReadUint32()
	// 	localAddrBytes := r.ReadBytes(4)
	// 	localPortBytes := r.ReadBytes(4)
	// 	remoteAddrBytes := r.ReadBytes(4)
	// 	remoteAddr := r.ReadIp(remoteAddrBytes)
	// 	remotePortBytes := r.ReadBytes(4)
	// 	remotePort := r.ReadPort(remotePortBytes)
	// 	pid := r.ReadUint32()
	// 	offLoadState := r.ReadUint32()

	// 	const (
	// 		BING string = "202.89.233.104" // netstat -ano | findstr 202.89.233.104
	// 		// kaola
	// 	)
	// 	testRemoteAddr := BING
	// 	if remoteAddr == "0.0.0.0" || remoteAddr == "127.0.0.1" || remoteAddr != testRemoteAddr {
	// 		continue
	// 	}
	// 	if pid > 0 && remoteAddr == testRemoteAddr {
	// 		row := &MIB_TCPROW2{DWORD(12), DWORD(bigEndian.Uint32(localAddrBytes)), DWORD(bigEndian.Uint32(localPortBytes)), DWORD(bigEndian.Uint32(remoteAddrBytes)), DWORD(bigEndian.Uint32(remotePortBytes)), DWORD(pid), TCP_CONNECTION_OFFLOAD_STATE(offLoadState)}
	// 		fmt.Println("row: ", row)

	// 		setTcpEntry := call.NewProc("SetTcpEntry")
	// 		r, _, _ := setTcpEntry.Call(uintptr(unsafe.Pointer(row)))
	// 		fmt.Println("setTcpEntry: ", r)
	// 	}
	// 	fmt.Printf("\t%-6d\t%s:%-16d\t%s:%-16d\t%d\t%d\n", state, r.ReadIp(localAddrBytes), r.ReadPort(localPortBytes), remoteAddr, remotePort, pid, offLoadState)
	// }
}
