package common

import "github.com/lxn/walk"

const (
	ICON_ADD = iota
	ICON_DELETE
	ICON_NEW
	ICON_OPEN
	ICON_PENCIL
	ICON_PLUS
	ICON_REFRESH
	ICON_SETTING
)

var IconMap = make(map[int]*walk.Bitmap)
var BrowserMap = make(map[string]bool)

func Init() {
	IconMap[ICON_ADD], _ = walk.NewBitmapFromFile("res/add.png")
	IconMap[ICON_DELETE], _ = walk.NewBitmapFromFile("res/delete.png")
	IconMap[ICON_NEW], _ = walk.NewBitmapFromFile("res/new.png")
	IconMap[ICON_OPEN], _ = walk.NewBitmapFromFile("res/open.png")
	IconMap[ICON_PENCIL], _ = walk.NewBitmapFromFile("res/pencil.png")
	IconMap[ICON_PLUS], _ = walk.NewBitmapFromFile("res/plus.png")
	IconMap[ICON_REFRESH], _ = walk.NewBitmapFromFile("res/refresh.png")
	IconMap[ICON_SETTING], _ = walk.NewBitmapFromFile("res/setting.png")

	BrowserMap["chrome.exe"] = true
	BrowserMap["firefox.exe"] = true
	BrowserMap["opera.exe"] = true
	BrowserMap["iexplore.exe"] = true
	BrowserMap["microsoftedge.exe"] = true
	BrowserMap["microsoftedgecp.exe"] = true
	BrowserMap["sogouexplorer.exe"] = true
	BrowserMap["qqbrowser.exe"] = true
	BrowserMap["360se.exe"] = true
	BrowserMap["360chrome.exe"] = true
	BrowserMap["liebao.exe"] = true
	BrowserMap["maxthon.exe"] = true
	BrowserMap["ucbrowser.exe"] = true
}
