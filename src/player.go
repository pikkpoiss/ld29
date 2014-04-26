package main

import (
	"../libs/twodee"
)

type Player struct {
	*twodee.BaseEntity
	Health      float64
	eventSystem *twodee.GameEventHandler
}

func (p *Player) MoveToListener() {}

func NewPlayer(e *twodee.BaseEntity, eventSystem *twodee.GameEventHandler) (player *Player) {
	player = &Player{
		e,
		100.0,
		eventSystem,
	}
	return
}
