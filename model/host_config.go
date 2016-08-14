package model

import (
	"github.com/lxn/walk"
	"github.com/tinycedar/lily/common"
)

type HostConfigModel struct {
	walk.TreeModelBase
	Roots []*HostConfigItem
}

func NewHostConfigModel() *HostConfigModel {
	//TODO laod system hosts for backup
	model := new(HostConfigModel)
	icon, err := walk.NewBitmapFromFile("res/new.png")
	if err != nil {
		common.Error("Error loading icon: ", err)
	}
	item := &HostConfigItem{Name: "test1", Icon: icon}
	model.Roots = append(model.Roots, item)
	model.Roots = append(model.Roots, &HostConfigItem{Name: "test2", Icon: icon})
	model.Roots = append(model.Roots, &HostConfigItem{Name: "test3", Icon: icon})
	model.Roots = append(model.Roots, &HostConfigItem{Name: "pre", Icon: icon})
	return model
}

func (m *HostConfigModel) Append(item *HostConfigItem) {
	m.Roots = append(m.Roots, item)
}

func (m *HostConfigModel) RootAt(index int) walk.TreeItem {
	return m.Roots[index]
}

func (m *HostConfigModel) RootCount() int {
	return len(m.Roots)
}

func (m *HostConfigModel) LazyPopulation() bool {
	return false
}

type HostConfigItem struct {
	walk.TreeItem
	Name string
	Icon *walk.Bitmap
}

var _ walk.TreeItem = new(HostConfigItem)

func (i *HostConfigItem) Image() interface{} {
	return i.Icon
}

func (i *HostConfigItem) Text() string {
	return i.Name
}

func (i *HostConfigItem) Parent() walk.TreeItem {
	return nil
}

func (i *HostConfigItem) ChildCount() int {
	return 0
}

func (i *HostConfigItem) ChildAt(index int) walk.TreeItem {
	return nil
}
