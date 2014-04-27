package main

import (
	"../libs/twodee"
)

type ItemType int

// Corresponds with frame indices
const (
	_                 = iota
	ItemDown ItemType = iota
	ItemUp
	Item1
	Item2
	Item3
	Item4
	ItemFinal
	item_sentinel
)

const (
	NumberOfItemTypes = int(item_sentinel)
)

type Item struct {
	*twodee.BaseEntity
	Id   ItemType
	Name string
}

func NewItem(itemType ItemType, name string, x, y, w, h float32) (item *Item) {
	item = &Item{
		BaseEntity: twodee.NewBaseEntity(x, y, w, h, 0, int(itemType)),
		Id:         itemType,
		Name:       name,
	}
	return
}

func (i *Item) getType() (itemType ItemType) {
	return i.Id
}
