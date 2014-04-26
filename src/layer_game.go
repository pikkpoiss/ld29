package main

import (
	twodee "../libs/twodee"
	"time"
)

type GameLayer struct {
	Level         *Level
	BatchRenderer *twodee.BatchRenderer
	Bounds        twodee.Rectangle
	App           *Application
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{
		App:    app,
		Bounds: twodee.Rect(0, 0, 10, 10),
	}
	if layer.BatchRenderer, err = twodee.NewBatchRenderer(layer.Bounds, app.WinBounds); err != nil {
		return
	}
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

func (l *GameLayer) Render() {
	var err error
	l.BatchRenderer.Bind()
	for _, geom := range l.Level.Geometry {
		if err = l.BatchRenderer.Draw(geom, 0, 0, 0); err != nil {
			panic(err)
		}
	}
	l.BatchRenderer.Unbind()
	return
}

func (l *GameLayer) Update(elapsed time.Duration) {
}

func (l *GameLayer) Reset() (err error) {
	return
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.KeyEscape:
		l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(GameIsClosing))
	}
	return true
}
