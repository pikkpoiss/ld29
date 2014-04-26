package main

import (
	"../libs/twodee"
)

type Player struct {
	*twodee.BaseEntity
	Health float64
}

func NewPlayer(e *twodee.BaseEntity) (player *Player) {
	player = &Player{
		e,
		100.0,
	}
	return
}
