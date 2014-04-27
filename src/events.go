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
	MenuPauseMusic
	MenuResumeMusic
	MenuMove
	MenuSelect
	PlayerTouchedItem
	PlayerUsedItem
	PlayerDestroyedItem
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
	twodee.BasicGameEvent
	Dir     MoveDirection
	Inverse bool
}

func NewPlayerMoveEvent(direction MoveDirection) (e *PlayerMoveEvent) {
	e = &PlayerMoveEvent{
		*twodee.NewBasicGameEvent(PlayerMove),
		direction,
		false,
	}
	return
}

func NewInversePlayerMoveEvent(direction MoveDirection) (e *PlayerMoveEvent) {
	e = &PlayerMoveEvent{
		*twodee.NewBasicGameEvent(PlayerMove),
		direction,
		true,
	}
	return
}

type PlayerTouchedItemEvent struct {
	twodee.BasicGameEvent
	Item *Item
}

func NewPlayerTouchedItemEvent(i *Item) (e *PlayerTouchedItemEvent) {
	e = &PlayerTouchedItemEvent{
		*twodee.NewBasicGameEvent(PlayerTouchedItem),
		i,
	}
	return
}

type PlayerUsedItemEvent struct {
	twodee.BasicGameEvent
	Item *Item
}

func NewPlayerUsedItemEvent(i *Item) (e *PlayerUsedItemEvent) {
	e = &PlayerUsedItemEvent{
		*twodee.NewBasicGameEvent(PlayerUsedItem),
		i,
	}
	return
}

type PlayerDestroyedItemEvent struct {
	twodee.BasicGameEvent
	Item *Item
}

func NewPlayerDestroyedItemEvent(i *Item) (e *PlayerDestroyedItemEvent) {
	e = &PlayerDestroyedItemEvent{
		*twodee.NewBasicGameEvent(PlayerDestroyedItem),
		i,
	}
	return
}
