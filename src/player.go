package main

import (
	"../libs/twodee"
)

type Player struct {
	*twodee.BaseEntity
	Health           float64
	eventSystem      *twodee.GameEventHandler
	moveToObserverId int
}

func (p *Player) MoveToObserver(e twodee.GETyper) {}

func (p *Player) Delete() {
	p.eventSystem.RemoveObserver(PlayerMove, p.moveToObserverId)
}

func NewPlayer(e *twodee.BaseEntity, eventSystem *twodee.GameEventHandler) (player *Player) {
	player = &Player{
		e,
		100.0,
		eventSystem,
		-1,
	}
	player.moveToObserverId = eventSystem.AddObserver(PlayerMove, player.MoveToObserver)
	return
}
