package core

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/conf"
)

const (
	systemHosts = "C:/Windows/System32/drivers/etc/hosts"
)

var batcher *Batcher

func FireHostsSwitch() {
	common.Info("============================== Fire hosts switch ==============================")
	if batcher != nil {
		batcher.Close()
	}
	doProcess()
	batcher = initSystemHostsWatcher()
	go startSystemHostsWatcher()
}

func initSystemHostsWatcher() *Batcher {
	batcher, err := New(time.Millisecond * 300)
	if err != nil {
		common.Error("Fail to initialize batcher")
	}
	if err = batcher.Add(systemHosts); err != nil {
		common.Error("Fail to add system hosts: %s", systemHosts)
	}
	return batcher
}

func startSystemHostsWatcher() {
	if batcher == nil {
		common.Error("Fail to start system hosts watcher, watcher is nil")
		return
	}
	for events := range batcher.Events {
		for _, event := range events {
			if event.Op&fsnotify.Write == fsnotify.Write {
				common.Info("modified file: %v", event)
				doProcess()
				break
			}
		}
	}
	// never return
}

// 1. Find collection of same domain names between system hosts and currentHostIndex
// 2. Disconnect the TCP connections(http:80 & https:443) of collection found above
func doProcess() {
	overlapHostConfigMap := getOverlapHostConfigMap()
	common.Info("overlapHostConfigMap: %v", overlapHostConfigMap)
	if len(overlapHostConfigMap) == 0 {
		return
	}
	table := getTCPTable()
	for i := uint32(0); i < uint32(table.dwNumEntries); i++ {
		row := table.table[i]
		if row.dwOwningPid <= 0 {
			continue
		}
		ip := row.displayIP(row.dwRemoteAddr)
		port := row.displayPort(row.dwRemotePort)
		if _, ok := overlapHostConfigMap[ip]; !ok {
			continue
		}
		if port == 80 || port == 443 {
			if err := CloseTCPEntry(row); err != nil {
				common.Error("Fail to close TCP connections: Pid = %v, Addr = %v:%v\n", row.dwOwningPid, ip, port)
			} else {
				common.Info("Succeed to close TCP connections: Pid = %v, Addr = %v:%v", row.dwOwningPid, ip, port)
			}
		}
	}
}

func getOverlapHostConfigMap() map[string]bool {
	result := make(map[string]bool)
	var currentHostConfigMap map[string]string
	if current := conf.Config.HostConfigModel.RootAt(conf.Config.CurrentHostIndex); current != nil {
		currentHostConfigMap = readHostConfigMap("conf/hosts/" + current.Text() + ".hosts")
	}
	common.Info("currentHostConfigMap: %v", currentHostConfigMap)
	if len(currentHostConfigMap) == 0 {
		return result
	}
	common.Info("systemConfigMap: %v", readHostConfigMap(systemHosts))
	for k, v := range readHostConfigMap(systemHosts) {
		v2, ok := currentHostConfigMap[k]
		common.Info("k: %s, %s - %s", k, v, v2)
		if ok && v != v2 {
			result[v] = true
		}
	}
	return result
}

func process() {
	// hostConfigMap := readFile()
	browserProcessMap := getBrowserProcessMap()
	table := getTCPTable()
	// group by process
	tcpRowByProcessNameMap := make(map[string][]*MIB_TCPROW2)
	for i := uint32(0); i < uint32(table.dwNumEntries); i++ {
		row := table.table[i]
		if row.dwOwningPid <= 0 {
			continue
		}
		// remoteAddr := row.displayIP(row.dwRemoteAddr)
		// if _, ok := hostConfigMap[remoteAddr]; ok {
		// common.Info("====== remoteAddr= %v\tbrowserProcessMap = %v\tpid = %v", remoteAddr, browserProcessMap, row.dwOwningPid)
		if processName, ok := browserProcessMap[uint32(row.dwOwningPid)]; ok {
			pidSlice, ok := tcpRowByProcessNameMap[processName]
			if !ok {
				pidSlice = []*MIB_TCPROW2{}
			}
			pidSlice = append(pidSlice, row)
			tcpRowByProcessNameMap[processName] = pidSlice
		}
		// }
		// common.Info("\t%-6d\t%s:%-16d\t%s:%-16d\t%d\t%d\n", row.dwState, row.displayIP(row.dwLocalAddr), row.displayPort(row.dwLocalPort), row.displayIP(row.dwRemoteAddr), row.displayPort(row.dwRemotePort), row.dwOwningPid, row.dwOffloadState)
	}
	browsers := []string{}
	for k := range tcpRowByProcessNameMap {
		browsers = append(browsers, k)
	}
	common.Info("Browsers: %v", browsers)
	for processName, rowSlice := range tcpRowByProcessNameMap {
		success := true
		for _, row := range rowSlice {
			if err := CloseTCPEntry(row); err != nil {
				success = false
				common.Error("Fail to close TCP connections: %s, %v\n", processName, row.dwOwningPid)
			}
		}
		if success {
			common.Info("Succeed to close TCP connections: %s", processName)
		}
	}
}

func readHostConfigMap(path string) map[string]string {
	hostConfigMap := make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		common.Error("Fail to open system_hosts: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			config := trimDuplicateSpaces(line)
			if len(config) == 2 {
				hostConfigMap[config[1]] = config[0]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		common.Error("Fail to read system_hosts: %s", err)
	}
	return hostConfigMap
}

func trimDuplicateSpaces(line string) []string {
	temp := []string{}
	line = strings.TrimSpace(line)
	for _, v := range strings.SplitN(line, " ", 2) {
		if trimed := strings.TrimSpace(v); trimed != "" {
			temp = append(temp, trimed)
		}
	}
	return temp
}
