package main

type ItemType int

const (
	Item1 ItemType = iota
	Item2
	Item3
	Item4
	Item5
	item_sentinel
)

const (
	NumberOfItemTypes = int(item_sentinel)
)

type Item struct {
	Id   ItemType
	Name string
}

func NewItem(itemType ItemType, name string) (item *Item) {
	item = &Item{
		Id:   itemType,
		Name: name,
	}
	return
}

func (i *Item) getType() (itemType ItemType) {
	return i.Id
}
