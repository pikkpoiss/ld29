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

const (
	North = iota
	East
	South
	West
)

type PlayerMoveEvent struct {
	*twodee.BasicGameEvent
	moveDir int
}

func NewPlayerMoveEvent(direction int) (e *PlayerMoveEvent) {
	e = &PlayerMoveEvent{
		twodee.NewBasicGameEvent(PlayerMove),
		direction,
	}
	return
}
