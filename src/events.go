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
	PlayExploreMusic
	PauseMusic
	ResumeMusic
	PlayerPickedUpItem
	sentinel
)

const (
	NumGameEventTypes = int(sentinel)
)

type MoveDirection int

const (
	North MoveDirection = iota
	East
	South
	West
	None
)

type PlayerMoveEvent struct {
	*twodee.BasicGameEvent
	Dir MoveDirection
}

func NewPlayerMoveEvent(direction MoveDirection) (e *PlayerMoveEvent) {
	e = &PlayerMoveEvent{
		twodee.NewBasicGameEvent(PlayerMove),
		direction,
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
