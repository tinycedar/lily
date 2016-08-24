package core

import (
	"github.com/tinycedar/lily/common"
	"golang.org/x/sys/windows/registry"
)

// OpenRegistry set DNS related registry
func OpenRegistry() {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		common.Error("Error Open registry key: %v", err)
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
