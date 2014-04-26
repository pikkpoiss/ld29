package main

import (
	twodee "../libs/twodee"
)

const (
	UpLayer twodee.GameEventType = iota
	DownLayer
	UpWaterLevel
	DownWaterLevel
	GameIsClosing
	sentinel
)

const (
	NumGameEventTypes = int(sentinel)
)
