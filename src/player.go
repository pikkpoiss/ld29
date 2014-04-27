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
	State       EntityState
	CanGetItem  bool
}

type EntityState int32

const (
	_                    = iota
	Standing EntityState = 1 << iota
	Walking
	Left
	Right
	Up
	Down
	ClimbUp
)

var PlayerAnimations = map[EntityState][]int{
	Standing | Up:    []int{24},
	Standing | Down:  []int{8},
	Standing | Left:  []int{16},
	Standing | Right: []int{16},
	Walking | Up:     []int{25, 26, 27, 28, 29, 30},
	Walking | Down:   []int{9, 10, 11, 12, 13, 14},
	Walking | Left:   []int{17, 18, 19, 20, 21, 22},
	Walking | Right:  []int{17, 18, 19, 20, 21, 22},
	ClimbUp | Down:   []int{32, 33, 34, 35, 36, 37, 38, 8},
}

func NewPlayer(x, y float32) (player *Player) {
	var (
		inv = make([]*Item, 0, NumberOfItemTypes)
	)
	player = &Player{
		AnimatingEntity: twodee.NewAnimatingEntity(
			x, y,
			32.0/PxPerUnit, 32.0/PxPerUnit,
			0,
			twodee.Step10Hz,
			[]int{8},
		),
		Health:      100.0,
		Speed:       PlayerBaseSpeed,
		Velocity:    twodee.Pt(0, 0),
		DesiredMove: None,
		Inventory:   inv,
		CanGetItem:  true,
	}
	return
}

func (p *Player) RemState(state EntityState) {
	p.SetState(p.State & ^state)
}

func (p *Player) AddState(state EntityState) {
	p.SetState(p.State | state)
}

func (p *Player) SwapState(rem, add EntityState) {
	p.SetState(p.State & ^rem | add)
}

func (p *Player) SetState(state EntityState) {
	if state != p.State {
		p.State = state
		if frames, ok := PlayerAnimations[p.State]; ok {
			p.SetFrames(frames)
		}
	}
}

func (p *Player) FlippedX() bool {
	return p.State&Left > 0
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
		p.SwapState(Walking, Standing)
		return
	case North:
		a = twodee.Pt(bounds.Min.X+Fudge, bounds.Max.Y+p.Speed)
		b = twodee.Pt(bounds.Max.X-Fudge, bounds.Max.Y+p.Speed)
		pos.Y += p.Speed
		trunc = l.GridAlignedY(l.Active, pos)
		p.SetState(Walking | Up)
	case South:
		a = twodee.Pt(bounds.Min.X+Fudge, bounds.Min.Y-p.Speed)
		b = twodee.Pt(bounds.Max.X-Fudge, bounds.Min.Y-p.Speed)
		pos.Y -= p.Speed
		trunc = l.GridAlignedY(l.Active, pos)
		p.SetState(Walking | Down)
	case East:
		a = twodee.Pt(bounds.Max.X+p.Speed, bounds.Min.Y+Fudge)
		b = twodee.Pt(bounds.Max.X+p.Speed, bounds.Max.Y-Fudge)
		pos.X += p.Speed
		trunc = l.GridAlignedX(l.Active, pos)
		p.SetState(Walking | Right)
	case West:
		a = twodee.Pt(bounds.Min.X-p.Speed, bounds.Min.Y+Fudge)
		b = twodee.Pt(bounds.Min.X-p.Speed, bounds.Max.Y-Fudge)
		pos.X -= p.Speed
		trunc = l.GridAlignedX(l.Active, pos)
		p.SetState(Walking | Left)
	}
	if l.FrontierCollides(l.Active, a, b) {
		p.MoveTo(trunc)
	} else {
		p.MoveTo(pos)
	}
}

const (
	PlayerBaseSpeed      = 0.2
	PlayerFastSpeed      = 0.3
	PlayerSuperFastSpeed = 0.4
)

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
		if p.Speed < PlayerFastSpeed {
			p.Speed = PlayerFastSpeed
		}
	case ItemFinal:
		if p.Speed < PlayerSuperFastSpeed {
			p.Speed = PlayerSuperFastSpeed
		}
	}
}
