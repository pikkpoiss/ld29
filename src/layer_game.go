package main

import (
	twodee "../libs/twodee"
	"time"
)

type GameLayer struct {
	Level           *Level
	BatchRenderer   *twodee.BatchRenderer
	Bounds          twodee.Rectangle
	App             *Application
	currentLayer    int
	layerTransition Tween
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{
		App:    app,
		Bounds: twodee.Rect(0, 0, 10, 10),
	}
	if layer.BatchRenderer, err = twodee.NewBatchRenderer(layer.Bounds, app.WinBounds); err != nil {
		return
	}
	layer.layerTransition = NewLinearTween(0, 10, time.Duration(5)*time.Second)
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
	var (
		err error
		y   float32 = 0.0
	)
	l.BatchRenderer.Bind()
	geom := l.Level.Geometry[l.currentLayer]
	if l.layerTransition != nil {
		y = l.layerTransition.Current()
	}
	if err = l.BatchRenderer.Draw(geom, 0, y, 0); err != nil {
		panic(err)
	}
	l.BatchRenderer.Unbind()
	return
}

func (l *GameLayer) Update(elapsed time.Duration) {
	if l.layerTransition != nil {
		if l.layerTransition.Done() {
			l.layerTransition = nil
		} else {
			l.layerTransition.Update(elapsed)
		}
	}
}

func (l *GameLayer) Reset() (err error) {
	return
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.KeyEvent:
		if event.Type == twodee.Release {
			break
		}
		switch event.Code {
		case twodee.KeyEscape:
			l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(GameIsClosing))
		case twodee.KeyLeft:
		case twodee.KeyRight:
		case twodee.KeyUp:
		case twodee.KeyDown:
		}
	}
	return true
}
