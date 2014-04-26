package main

import (
	twodee "../libs/twodee"
)

const (
	UpLayer twodee.GameEventType = iota
	DownLayer
	UpWaterLevel
	DownWaterLevel
	PlayerMoveEventType
	sentinel
)

const (
	NumGameEventTypes = int(sentinel)
)

type PlayerMoveEvent struct {
	*twodee.BasicGameEvent
	EndPoint *twodee.Point
}

func NewPlayerMoveEvent(pt *twodee.Point) (e *PlayerMoveEvent) {
	e = &PlayerMoveEvent{
		twodee.NewBasicGameEvent(PlayerMoveEventType),
		EndPoint: pt,
	}
	return
}
