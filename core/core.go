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

func doProcess() bool {
	success := true
	hostsIPMap := getHostsIpMap()
	overwriteSystemHosts()
	processNameMap := GetProcessNameMap()
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
		process := strings.ToLower(processNameMap[uint32(row.dwOwningPid)])
		if hostsIPMap[ip] || common.BrowserMap[process] {
			if err := CloseTCPEntry(row); err != nil {
				common.Error("Fail to close TCP connections: Process = %v, Pid = %v, Addr = %v:%v", process, row.dwOwningPid, ip, port)
				success = false
			} else {
				common.Info("Succeed to close TCP connections: Process = %v, Pid = %v, Addr = %v:%v", process, row.dwOwningPid, ip, port)
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
