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
		BaseEntity:  e,
		Health:      100.0,
		Speed:       0.2,
		Velocity:    twodee.Pt(0, 0),
		DesiredMove: None,
	}
	return
}

const Fudge = 0.01

func (p *Player) AttemptMove(l *Level) {
	var (
		a, b, trunc twodee.Point
		bounds      = p.Bounds()
		pos         = p.Pos()
	)
	switch p.DesiredMove {
	case None:
		return
	case North:
		a = twodee.Pt(bounds.Min.X+Fudge, bounds.Max.Y+p.Speed)
		b = twodee.Pt(bounds.Max.X-Fudge, bounds.Max.Y+p.Speed)
		pos.Y += p.Speed
		trunc = l.GridAlignedY(l.Active, pos)
	case South:
		a = twodee.Pt(bounds.Min.X+Fudge, bounds.Min.Y-p.Speed)
		b = twodee.Pt(bounds.Max.X-Fudge, bounds.Min.Y-p.Speed)
		pos.Y -= p.Speed
		trunc = l.GridAlignedY(l.Active, pos)
	case East:
		a = twodee.Pt(bounds.Max.X+p.Speed, bounds.Min.Y+Fudge)
		b = twodee.Pt(bounds.Max.X+p.Speed, bounds.Max.Y-Fudge)
		pos.X += p.Speed
		trunc = l.GridAlignedX(l.Active, pos)
	case West:
		a = twodee.Pt(bounds.Min.X-p.Speed, bounds.Min.Y+Fudge)
		b = twodee.Pt(bounds.Min.X-p.Speed, bounds.Max.Y-Fudge)
		pos.X -= p.Speed
		trunc = l.GridAlignedX(l.Active, pos)
	}
	if l.FrontierCollides(l.Active, a, b) {
		p.MoveTo(trunc)
	} else {
		p.MoveTo(pos)
	}
}
