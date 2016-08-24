package core

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tinycedar/lily/common"
)

const (
	system_hosts = "C:/Windows/System32/drivers/etc/hosts"
)

func NewSystemHostsWatcher() {
	batcher, err := New(time.Millisecond * 300)
	if err == nil {
		if err = batcher.Add(system_hosts); err != nil {
			log.Fatal(err)
		}
		for events := range batcher.Events {
			// fmt.Println("events: ", events)
			for _, event := range events {
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event)
					process()
					break
				}
			}
		}
	}
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
		// fmt.Println("====== remoteAddr= ", remoteAddr, "\tbrowserProcessMap = ", browserProcessMap, "\tpid = ", row.dwOwningPid)
		if processName, ok := browserProcessMap[uint32(row.dwOwningPid)]; ok {
			pidSlice, ok := tcpRowByProcessNameMap[processName]
			if !ok {
				pidSlice = []*MIB_TCPROW2{}
			}
			pidSlice = append(pidSlice, row)
			tcpRowByProcessNameMap[processName] = pidSlice
		}
		// }
		// fmt.Printf("\t%-6d\t%s:%-16d\t%s:%-16d\t%d\t%d\n", row.dwState, row.displayIP(row.dwLocalAddr), row.displayPort(row.dwLocalPort), row.displayIP(row.dwRemoteAddr), row.displayPort(row.dwRemotePort), row.dwOwningPid, row.dwOffloadState)
	}
	common.Info("==================== Running browser =====================")
	for k := range tcpRowByProcessNameMap {
		common.Info("%v", k)
	}
	common.Info("================== Execute Result  =====================")
	for processName, rowSlice := range tcpRowByProcessNameMap {
		success := true
		for _, row := range rowSlice {
			if err := CloseTCPEntry(row); err != nil {
				success = false
				common.Error("Fail to close TCP connections: %s, %v\n", processName, row.dwOwningPid)
			}
		}
		if success {
			common.Info("Succeed to close TCP connections: %s\n", processName)
		}
	}
}

// func readFile() map[string]string {
// 	hostConfigMap := make(map[string]string)
// 	file, err := os.Open(system_hosts)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()
// 	scanner := bufio.NewScanner(file)
// 	// fmt.Println("============================== Reading file begin =====================================")
// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())
// 		if line != "" && !strings.HasPrefix(line, "#") {
// 			config := strings.Split(scanner.Text(), " ")
// 			if len(config) != 2 {
// 				config = strings.Split(scanner.Text(), "\t")
// 			}
// 			if len(config) == 2 {
// 				// fmt.Println(config[1], "\t", config[0])
// 				hostConfigMap[config[0]] = config[1]
// 			}
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	// fmt.Println("============================== Reading file end =====================================")
// 	return hostConfigMap
// }
