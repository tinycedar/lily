package core

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tinycedar/lily/common"
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
	process()
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
				process()
				break
			}
		}
	}
	// never return
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

// func readFile() map[string]string {
// 	hostConfigMap := make(map[string]string)
// 	file, err := os.Open(system_hosts)
// 	if err != nil {
// 		common.Error("Fail to open system_hosts: %s", err)
// 	}
// 	defer file.Close()
// 	scanner := bufio.NewScanner(file)
// 	// common.Info("============================== Reading file begin =====================================")
// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())
// 		if line != "" && !strings.HasPrefix(line, "#") {
// 			config := strings.Split(scanner.Text(), " ")
// 			if len(config) != 2 {
// 				config = strings.Split(scanner.Text(), "\t")
// 			}
// 			if len(config) == 2 {
// 				// common.Info("%v\t%v", config[1], config[0])
// 				hostConfigMap[config[0]] = config[1]
// 			}
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		common.Error("Fail to read system_hosts: %s", err)
// 	}
// 	// common.Info("============================== Reading file end =====================================")
// 	return hostConfigMap
// }
