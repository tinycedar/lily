package core

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/tinycedar/lily/common"
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
	DWORD                        uint32
	TCP_CONNECTION_OFFLOAD_STATE uint32
)

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
	return binary.BigEndian.Uint16([]byte{byte(val), byte(val >> 8)})
}

func newTCPRow(r *ClassReader) *MIB_TCPROW2 {
	return &MIB_TCPROW2{DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), DWORD(r.ReadUint32()), TCP_CONNECTION_OFFLOAD_STATE(r.ReadUint32())}
}

type MIB_TCPTABLE2 struct {
	dwNumEntries DWORD
	table        []*MIB_TCPROW2
}

func (t *MIB_TCPTABLE2) String() string {
	common.Info("================  tcp table ======================= %v", t.dwNumEntries)
	for i := uint32(0); i < uint32(t.dwNumEntries); i++ {
		row := t.table[i]
		common.Info("%v\t%v:%v", row, row.displayIP(row.dwRemoteAddr), row.displayPort(row.dwRemotePort))
	}
	common.Info("================  tcp table end =======================")
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

func getTCPTable() *MIB_TCPTABLE2 {
	getTCPTable2 := syscall.NewLazyDLL("Iphlpapi.dll").NewProc("GetTcpTable2")
	var n uint32
	if err, _, _ := getTCPTable2.Call(uintptr(unsafe.Pointer(&MIB_TCPTABLE2{})), uintptr(unsafe.Pointer(&n)), 1); syscall.Errno(err) != syscall.ERROR_INSUFFICIENT_BUFFER {
		common.Error("Error calling GetTcpTable2: %v\n", syscall.Errno(err))
	}
	b := make([]byte, n)
	if err, _, _ := getTCPTable2.Call(uintptr(unsafe.Pointer(&b[0])), uintptr(unsafe.Pointer(&n)), 1); err != 0 {
		common.Error("Error calling GetTcpTable2: %v\n", syscall.Errno(err))
	}
	const (
		// netstat -ano | findstr 202.89.233.104
		LOCALHOST string = "127.0.0.1"
		BING      string = "202.89.233.103"
		KAOLA     string = "127.0.0.1"
	)
	table := newTCPTable(NewClassReader(b))
	return table
}

func CloseTCPEntry(row *MIB_TCPROW2) error {
	row.dwState = 12
	if err, _, _ := syscall.NewLazyDLL("Iphlpapi.dll").NewProc("SetTcpEntry").Call(uintptr(unsafe.Pointer(row))); err != 0 {
		return syscall.Errno(err)
	}
	return nil
}
