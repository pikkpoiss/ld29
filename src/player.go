package main

import (
	"../libs/twodee"
)

type Player struct {
	*twodee.BaseEntity
	Health      float32
	Speed       float32
	Velocity    twodee.Point
	DesiredMove MoveDirection
}

func NewPlayer(e *twodee.BaseEntity) (player *Player) {
	player = &Player{
		BaseEntity: e,
		Health:     100.0,
		Speed:      0.1,
		Velocity:   twodee.Pt(0, 0),
	}
	return
}

func (p *Player) AttemptMove(l *Level) {
}
