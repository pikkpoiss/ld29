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
	playerLayer   int
	topLayer      int
	topTrans      Tween
	botLayer      int
	botTrans      Tween
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
	l.LayerAdvance()
	return
}

func (l *GameLayer) Delete() {
	if l.Level != nil {
		l.Level.Delete()
	}
}

func (l *GameLayer) Render() {
	var (
		y           float32 = 0.0
		drawnPlayer bool    = false
	)
	l.BatchRenderer.Bind()
	for i := l.Level.Layers - 1; i >= 0; i-- {
		if i == l.botLayer && l.botTrans != nil {
			y = l.botTrans.Current()
			l.BatchRenderer.Draw(l.Level.Geometry[l.botLayer], 0, y, 0)
			if i == l.playerLayer {
				drawnPlayer = true
			}
		}
		if i == l.playerLayer && drawnPlayer == false {
			l.BatchRenderer.Draw(l.Level.Geometry[l.playerLayer], 0, 0, 0)
		}
		if i == l.topLayer && l.topTrans != nil {
			y = l.topTrans.Current()
			l.BatchRenderer.Draw(l.Level.Geometry[l.topLayer], 0, y, 0)
		}
	}
	l.BatchRenderer.Unbind()
	return
}

func (l *GameLayer) Update(elapsed time.Duration) {
	if l.topTrans != nil {
		if l.topTrans.Done() {
			l.topTrans = nil
		} else {
			l.topTrans.Update(elapsed)
		}
	}
	if l.botTrans != nil {
		if l.botTrans.Done() {
			l.botTrans = nil
		} else {
			l.botTrans.Update(elapsed)
		}
	}
}

func (l *GameLayer) LayerAdvance() {
	if l.playerLayer >= l.Level.Layers {
		return
	}
	l.topLayer = l.playerLayer
	l.botLayer = l.topLayer + 1
	l.playerLayer = l.botLayer
	l.topTrans = NewLinearTween(0, l.Bounds.Max.Y, time.Duration(500)*time.Millisecond)
	l.botTrans = NewLinearTween(-1, 0, time.Duration(200)*time.Millisecond)
}

func (l *GameLayer) LayerRewind() {
	if l.playerLayer <= 0 {
		return
	}
	l.topLayer = l.playerLayer - 1
	l.botLayer = l.topLayer + 1
	l.playerLayer = l.topLayer
	l.topTrans = NewLinearTween(l.Bounds.Max.Y, 0, time.Duration(500)*time.Millisecond)
	l.botTrans = NewLinearTween(0, -l.Bounds.Max.Y, time.Duration(200)*time.Millisecond)
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
