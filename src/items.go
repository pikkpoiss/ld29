package main

type ItemType int

const (
	item1 ItemType = iota
	item2
	item3
	item4
	item5
	item_sentinel
)

const (
	NumberOfItemTypes = int(item_sentinel)
)

type Item struct {
	Id   int
	Name string
}
