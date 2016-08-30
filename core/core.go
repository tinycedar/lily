package core

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/conf"
)

const (
	systemHosts = "C:/Windows/System32/drivers/etc/hosts"
)

// var batcher *Batcher

func FireHostsSwitch() bool {
	common.Info("============================== Fire hosts switch ==============================")
	// if batcher != nil {
	// 	batcher.Close()
	// }
	return doProcess()
	// batcher = initSystemHostsWatcher()
	// go startSystemHostsWatcher()
}

// 1. Find collection of same domain names between system hosts and currentHostIndex
// 2. Disconnect the TCP connections(http:80 & https:443) of collection found above
func doProcess() bool {
	success := true
	hostsIpMap := getHostsIpMap()
	overwriteSystemHosts()
	table := getTCPTable()
	for i := uint32(0); i < uint32(table.dwNumEntries); i++ {
		row := table.table[i]
		ip := row.displayIP(row.dwRemoteAddr)
		port := row.displayPort(row.dwRemotePort)
		if row.dwOwningPid <= 0 {
			continue
		}
		if port != 80 && port != 443 {
			continue
		}
		supported := common.BrowserMap[strings.ToLower(OpenProcess(uint32(row.dwOwningPid)))]
		if hostsIpMap[ip] || supported {
			if err := CloseTCPEntry(row); err != nil {
				common.Error("Fail to close TCP connections: Pid = %v, Addr = %v:%v\n", row.dwOwningPid, ip, port)
				success = false
			} else {
				common.Info("Succeed to close TCP connections: Pid = %v, Addr = %v:%v", row.dwOwningPid, ip, port)
			}
		}
	}
	return success
}

// find all the ip of system and current hosts
func getHostsIpMap() map[string]bool {
	ipMap := make(map[string]bool)
	for _, v := range readHostConfigMap(systemHosts) {
		ipMap[v] = true
	}
	model := conf.Config.HostConfigModel
	index := conf.Config.CurrentHostIndex
	if length := len(model.Roots); index >= 0 && length > 0 && index < length {
		path := "conf/hosts/" + model.RootAt(index).Text() + ".hosts"
		for _, v := range readHostConfigMap(path) {
			ipMap[v] = true
		}
	}
	return ipMap
}

func overwriteSystemHosts() {
	bytes := ReadCurrentHostConfig()
	if bytes == nil {
		return
	}
	if err := ioutil.WriteFile(systemHosts, bytes, os.ModeExclusive); err != nil {
		common.Error("Error writing to system hosts file: ", err)
	}
}

func ReadCurrentHostConfig() []byte {
	model := conf.Config.HostConfigModel
	index := conf.Config.CurrentHostIndex
	if length := len(model.Roots); index < 0 || length <= 0 || index >= length {
		return nil
	}
	bytes, err := ioutil.ReadFile("conf/hosts/" + model.RootAt(index).Text() + ".hosts")
	if err != nil {
		common.Info("Error reading host config: %v", err)
		return nil
	}
	return bytes
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
		common.Info("Fail to open system_hosts: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			config := strings.Fields(line)
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

// func initSystemHostsWatcher() *Batcher {
// 	batcher, err := New(time.Millisecond * 300)
// 	if err != nil {
// 		common.Error("Fail to initialize batcher")
// 	}
// 	if err = batcher.Add(systemHosts); err != nil {
// 		common.Error("Fail to add system hosts: %s", systemHosts)
// 	}
// 	return batcher
// }

// func startSystemHostsWatcher() {
// 	if batcher == nil {
// 		common.Error("Fail to start system hosts watcher, watcher is nil")
// 		return
// 	}
// 	for events := range batcher.Events {
// 		for _, event := range events {
// 			if event.Op&fsnotify.Write == fsnotify.Write {
// 				common.Info("modified file: %v", event)
// 				doProcess()
// 				break
// 			}
// 		}
// 	}
// 	// never return
// }
