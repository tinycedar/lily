package model

import "github.com/lxn/walk"

type HostConfigModel struct {
	walk.TreeModelBase
	Roots []*HostConfigItem
}

func NewHostConfigModel() *HostConfigModel {
	model := new(HostConfigModel)
	model.Roots = []*HostConfigItem{}
	return model
}

func (m *HostConfigModel) Append(item *HostConfigItem) {
	m.Roots = append(m.Roots, item)
	m.PublishItemsReset(nil)
}

func (m *HostConfigModel) Remove(item *HostConfigItem) bool {
	for i, size := 0, len(m.Roots); i < size; i++ {
		if item == m.Roots[i] {
			if size == 1 {
				m.Roots = nil
				m.PublishItemsReset(nil)
				return true
			}
			if i == 0 {
				m.Roots = m.Roots[1:]
			} else if i == size-1 {
				m.Roots = m.Roots[0 : size-1]
			} else {
				tmp := []*HostConfigItem{}
				tmp = append(tmp, m.Roots[:i]...)
				tmp = append(tmp, m.Roots[i+1:]...)
				m.Roots = tmp
			}
			m.PublishItemsReset(nil)
			return true
		}
	}
	return false
}

func (m *HostConfigModel) RemoveAll() {
	m.Roots = nil
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
