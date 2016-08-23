package main

import (
	"log"

	"golang.org/x/sys/windows/registry"

	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/core"
	"github.com/tinycedar/lily/gui"
)

// https://github.com/spf13/hugo/blob/master/watcher/batcher.go
func main() {
	go initBgProcessor()
	gui.InitMainWindow()
}

func initBgProcessor() {
	openRegistry()
	core.NewSystemHostsWatcher()
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

// func readFile() map[string]string {
// 	hostConfigMap := make(map[string]string)
// 	file, err := os.Open(hostFile)
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
