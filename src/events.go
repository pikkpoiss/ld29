package main

import (
	twodee "../libs/twodee"
)

const (
	UpLayer twodee.GameEventType = iota
	DownLayer
	sentinel
)

const (
	NumGameEventTypes = int(sentinel)
)
