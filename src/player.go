package main

import (
	"time"

	"../libs/twodee"
)

type DirectionsHistoryEntry struct {
	prev, next *DirectionsHistoryEntry
	dir        MoveDirection
}

type DirectionsHistory struct {
	tail       *DirectionsHistoryEntry
	directions []*DirectionsHistoryEntry
}

// Adds a MoveDirection to the history if it's not already present.
func (dh *DirectionsHistory) Add(d MoveDirection) {
	entry := dh.directions[d]
	if entry.prev == nil && entry.next == nil && entry != dh.tail {
		if dh.tail != nil {
			dh.tail.next = entry
			entry.prev = dh.tail
		}
		dh.tail = entry
	}
}

// Removes a particular MoveDirection from the history chain; sets its prev and next fields to nil. Resets dh's tail if necessary.
func (dh *DirectionsHistory) Remove(d MoveDirection) {
	entry := dh.directions[d]
	if entry.prev != nil {
		entry.prev.next = entry.next
	}
	if entry.next != nil {
		entry.next.prev = entry.prev
	}
	if entry == dh.tail {
		dh.tail = entry.prev
	}
	entry.prev = nil
	entry.next = nil
}

func (dh *DirectionsHistory) LatestDirection() (d MoveDirection) {
	if dh.tail != nil {
		return dh.tail.dir
	}
	return None
}

func NewDirectionsHistory() (dh *DirectionsHistory) {
	// Ugh, West+1 because we're generating a sparse array indexed on the
	// int value of directions; West should be the last one enumerated.
	dh = &DirectionsHistory{
		tail:       nil,
		directions: make([]*DirectionsHistoryEntry, West+1),
	}
	dirs := []MoveDirection{North, East, South, West}
	for _, d := range dirs {
		dh.directions[d] = &DirectionsHistoryEntry{
			prev: nil,
			next: nil,
			dir:  d,
		}
	}
	return
}

type Player struct {
	*twodee.AnimatingEntity
	EventSystem       *twodee.GameEventHandler
	MaxHealth         int32
	Health            int32
	Speed             float32
	Velocity          twodee.Point
	DirectionsHistory *DirectionsHistory
	DesiredMove       MoveDirection
	Inventory         []*Item
	State             EntityState
	CanGetItem        bool
	CanMove           bool
	IsPumping         bool
	destroyableItems  map[ItemId]bool
	HasFinalItem      bool
	LastUsed          *Item
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
	ClimbDown
)

const (
	Fudge                = 0.4
	PlayerBaseSpeed      = 0.2
	PlayerFastSpeed      = 0.3
	PlayerSuperFastSpeed = 0.4
	PlayerBaseHealth     = 1000
	PlayerWaterDamage    = 6
	PlayerHealthRegen    = 4
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
	ClimbDown | Down: []int{48, 49, 50, 51, 52, 53, 54, 8},
}

func NewPlayer(x, y float32, eventSystem *twodee.GameEventHandler) (player *Player) {
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
		EventSystem:       eventSystem,
		MaxHealth:         PlayerBaseHealth,
		Health:            PlayerBaseHealth,
		Speed:             PlayerBaseSpeed,
		Velocity:          twodee.Pt(0, 0),
		DirectionsHistory: NewDirectionsHistory(),
		DesiredMove:       None,
		Inventory:         inv,
		CanGetItem:        true,
		CanMove:           true,
		IsPumping:         false,
		HasFinalItem:      false,
		destroyableItems:  make(map[ItemId]bool),
		LastUsed:          nil,
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

// Updates the Player's desired movement direction as well as the affiliated data
// structures. If `invert`, then the movement key has been released and we should
// remove it from the affiliated data structures and perhaps pick a new movement
// direction from the tail of OrderedDirections.
func (p *Player) UpdateDesiredMove(d MoveDirection, invert bool) {
	if invert {
		p.DirectionsHistory.Remove(d)
		p.DesiredMove = p.DirectionsHistory.LatestDirection()
		return
	}
	// If the player is already moving in this direction, do nothing.
	if p.DesiredMove == d {
		return
	}
	p.DirectionsHistory.Add(d)
	p.DesiredMove = p.DirectionsHistory.LatestDirection()
}

func (p *Player) AttemptMove(l *Level) {
	if !p.CanMove {
		return
	}
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

func (p *Player) AddToInventory(item *Item) {
	p.EventSystem.Enqueue(twodee.NewBasicGameEvent(PlayPickupItemEffect))
	p.Inventory = append(p.Inventory, item)
	switch item.Id {
	case Item1:
		p.MaxHealth += 100
		p.Health += 100
	case Item2:
		p.MaxHealth += 100
		p.Health += 100
	case Item3:
		p.MaxHealth += PlayerBaseHealth
		p.Health += PlayerBaseHealth
	case Item4:
		if p.Speed < PlayerFastSpeed {
			p.Speed = PlayerFastSpeed
		}
	case ItemFinal:
		if p.Speed < PlayerSuperFastSpeed {
			p.Speed = PlayerSuperFastSpeed
		}
		p.HasFinalItem = true
	case ItemPickaxe:
		p.destroyableItems[ItemRock] = true
	}
}

func (p *Player) CanDestroy(item *Item) bool {
	return p.destroyableItems[item.Id]
}

func (p *Player) Update(elapsed time.Duration) {
	p.AnimatingEntity.Update(elapsed)
	if p.Health < p.MaxHealth {
		p.Health += PlayerHealthRegen
		if p.Health > p.MaxHealth {
			p.Health = p.MaxHealth
		}
	}
}

func (p *Player) Damage(damage int32) {
	p.Health -= damage
	if p.Health < 0 {
		p.Health = 0
	}
}

func (p *Player) HealthPercent() float32 {
	return float32(p.Health) / float32(p.MaxHealth)
}
