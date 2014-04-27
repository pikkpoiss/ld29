package main

import (
	"../libs/twodee"
)

type ItemType int
type ItemId int

// Corresponds with frame indices
const (
	_               = iota
	ItemDown ItemId = iota
	ItemUp
	Item1
	Item2
	Item3
	Item4
	ItemFinal
	item_sentinel
)

// Various item types
const (
	LayerThresholdItem ItemType = iota
	InventoryItem
	DestructableItem
)

var ItemIdToType = map[ItemId]ItemType{
	ItemDown:  LayerThresholdItem,
	ItemUp:    LayerThresholdItem,
	Item1:     InventoryItem,
	Item2:     InventoryItem,
	Item3:     InventoryItem,
	Item4:     InventoryItem,
	ItemFinal: InventoryItem,
}

const (
	NumberOfItemTypes = int(item_sentinel)
)

type Item struct {
	*twodee.BaseEntity
	Id   ItemId
	Type ItemType
	Name string
}

func NewItem(itemId ItemId, itemType ItemType, name string, x, y, w, h float32) (item *Item) {
	item = &Item{
		BaseEntity: twodee.NewBaseEntity(x, y, w, h, 0, int(itemId)),
		Id:         itemId,
		Type:       itemType,
		Name:       name,
	}
	return
}

func (i *Item) getType() ItemType {
	return i.Type
}
