package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"

	"github.com/fsnotify/fsnotify"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/gui"
)

const (
	hostFile = "C:/Windows/System32/drivers/etc/hosts"
)

// https://github.com/spf13/hugo/blob/master/watcher/batcher.go
func main() {
	go initBgProcessor()
	gui.InitMainWindow()
}

func initBgProcessor() {
	openRegistry()
	batcher, err := New(time.Millisecond * 300)
	if err == nil {
		if err = batcher.Add(hostFile); err != nil {
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

func openRegistry() {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()
	k.SetDWordValue("DnsCacheEnabled", 0x1)
	k.SetDWordValue("DnsCacheTimeout", 0x1)
	k.SetDWordValue("ServerInfoTimeOut", 0x1)
	if s, _, err := k.GetIntegerValue("DnsCacheEnabled"); err != nil {
		common.Error("Fail to get registry: DnsCacheEnabled", err)
	} else {
		common.Info("DnsCacheEnabled is %q\n", s)
	}

	if s, _, err := k.GetIntegerValue("DnsCacheTimeout"); err != nil {
		common.Error("Fail to get registry: DnsCacheTimeout", err)
	} else {
		common.Info("DnsCacheTimeout is %q\n", s)
	}

	if s, _, err := k.GetIntegerValue("ServerInfoTimeOut"); err != nil {
		common.Error("Fail to get registry: ServerInfoTimeOut", err)
	} else {
		common.Info("ServerInfoTimeOut is %q\n", s)
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

type Batcher struct {
	*fsnotify.Watcher
	interval time.Duration
	done     chan struct{}

	Events chan []fsnotify.Event // Events are returned on this channel
}

func New(interval time.Duration) (*Batcher, error) {
	watcher, err := fsnotify.NewWatcher()

	batcher := &Batcher{}
	batcher.Watcher = watcher
	batcher.interval = interval
	batcher.done = make(chan struct{}, 1)
	batcher.Events = make(chan []fsnotify.Event, 1)

	if err == nil {
		go batcher.run()
	}

	return batcher, err
}

func (b *Batcher) run() {
	tick := time.Tick(b.interval)
	evs := make([]fsnotify.Event, 0)
OuterLoop:
	for {
		select {
		case ev := <-b.Watcher.Events:
			evs = append(evs, ev)
		case <-tick:
			if len(evs) == 0 {
				continue
			}
			b.Events <- evs
			evs = make([]fsnotify.Event, 0)
		case <-b.done:
			break OuterLoop
		}
	}
	close(b.done)
}

func (b *Batcher) Close() {
	b.done <- struct{}{}
	b.Watcher.Close()
}

func readFile() map[string]string {
	hostConfigMap := make(map[string]string)
	file, err := os.Open(hostFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// fmt.Println("============================== Reading file begin =====================================")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			config := strings.Split(scanner.Text(), " ")
			if len(config) != 2 {
				config = strings.Split(scanner.Text(), "\t")
			}
			if len(config) == 2 {
				// fmt.Println(config[1], "\t", config[0])
				hostConfigMap[config[0]] = config[1]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// fmt.Println("============================== Reading file end =====================================")
	return hostConfigMap
}
