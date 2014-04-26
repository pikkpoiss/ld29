package main

import (
	"../libs/twodee"
)

type Item struct {
	*twodee.BaseEntity
	Id   int
	Name string
}
