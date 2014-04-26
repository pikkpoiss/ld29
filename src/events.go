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
