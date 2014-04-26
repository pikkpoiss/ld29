package main

import (
	"../libs/twodee"
)

type Player struct {
	*twodee.AnimatingEntity
	Health      float32
	Speed       float32
	Velocity    twodee.Point
	DesiredMove MoveDirection
	Inventory   []*Item
}

func NewPlayer(x, y float32) (player *Player) {
	var (
		inv = make([]*Item, 0, NumberOfItemTypes)
	)
	player = &Player{
		AnimatingEntity: twodee.NewAnimatingEntity(
			x, y,
			1, 1,
			0,
			twodee.Step10Hz,
			[]int{8},
		),
		Health:      100.0,
		Speed:       0.2,
		Velocity:    twodee.Pt(0, 0),
		DesiredMove: None,
		Inventory:   inv,
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

func (p *Player) AddToInventory(item *Item) {
	p.Inventory = append(p.Inventory, item)
	switch item.getType() {
	case Item1:
		p.Health = p.Health + 10
	case Item2:
		p.Health = p.Health + 20
	case Item3:
		p.Health = p.Health + 30
	case Item4:
		p.Speed = p.Speed + 10
	case Item5:
		p.Speed = p.Speed + 20
	}
}
