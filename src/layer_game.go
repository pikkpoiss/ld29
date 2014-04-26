package main

import ()

type GameLayer struct {
	Level *Level
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{}
	return
}

func (l *GameLayer) LoadLevel(path string) (err error) {
	l.Level, err = LoadLevel(path)
	return
}

func (l *GameLayer) Delete() {
	if l.Level != nil {
		l.Level.Delete()
	}
}
