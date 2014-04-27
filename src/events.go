package main

import (
	twodee "../libs/twodee"
)

const (
	UpLayer twodee.GameEventType = iota
	DownLayer
	UpWaterLevel
	DownWaterLevel
	PlayerMove
	GameIsClosing
	PlayOutdoorMusic
	PlayExploreMusic
	PlayWarningMusic
	PlayDangerMusic
	PauseMusic
	ResumeMusic
	PlayerPickedUpItem
	sentinel
)

const (
	NumGameEventTypes = int(sentinel)
)

type MoveDirection byte

const (
	None  MoveDirection = iota
	North MoveDirection = 1 << (iota - 1)
	East
	South
	West
)

type PlayerMoveEvent struct {
	*twodee.BasicGameEvent
	Dir     MoveDirection
	Inverse bool
}

func NewPlayerMoveEvent(direction MoveDirection) (e *PlayerMoveEvent) {
	e = &PlayerMoveEvent{
		twodee.NewBasicGameEvent(PlayerMove),
		direction,
		false,
	}
	return
}

func NewInversePlayerMoveEvent(direction MoveDirection) (e *PlayerMoveEvent) {
	e = &PlayerMoveEvent{
		twodee.NewBasicGameEvent(PlayerMove),
		direction,
		true,
	}
	return
}

type PlayerPickedUpItemEvent struct {
	*twodee.BasicGameEvent
	Item *Item
}

func NewPlayerPickedUpItemEvent(i *Item) (e *PlayerPickedUpItemEvent) {
	e = &PlayerPickedUpItemEvent{
		twodee.NewBasicGameEvent(PlayerPickedUpItem),
		i,
	}
	return
}
