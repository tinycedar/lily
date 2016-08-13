package model

import (
	"github.com/lxn/walk"
	"github.com/tinycedar/lily/common"
)

type hostConfigModel struct {
	walk.TreeModelBase
	roots []*hostConfigItem
}

func NewHostConfigModel() (*hostConfigModel, error) {
	model := new(hostConfigModel)
	icon, err := walk.NewBitmapFromFile("res/open.png")
	if err != nil {
		common.Error("Error loading icon: ", err)
	}
	model.roots = append(model.roots, &hostConfigItem{name: "Daniel", icon: icon})
	model.roots = append(model.roots, &hostConfigItem{name: "Elim", icon: icon})
	model.roots = append(model.roots, &hostConfigItem{name: "大家好", icon: icon})
	return model, nil
}

func (m *hostConfigModel) RootAt(index int) walk.TreeItem {
	return m.roots[index]
}

func (m *hostConfigModel) RootCount() int {
	return len(m.roots)
}

func (m *hostConfigModel) LazyPopulation() bool {
	return false
}

type hostConfigItem struct {
	walk.TreeItem
	name string
	icon *walk.Bitmap
}

var _ walk.TreeItem = new(hostConfigItem)

func (i *hostConfigItem) Image() interface{} {
	return i.icon
}

func (i *hostConfigItem) Text() string {
	return i.name
}

func (i *hostConfigItem) Parent() walk.TreeItem {
	return nil
}

func (i *hostConfigItem) ChildCount() int {
	return 0
}

func (i *hostConfigItem) ChildAt(index int) walk.TreeItem {
	return nil
}
