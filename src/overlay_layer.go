package main

import (
	"../libs/twodee"
	"time"
)

type OverlayLayer struct {

	Game              *GameLayer
	Events            *twodee.GameEventHandler
	TileRenderer      *twodee.TileRenderer
	Bounds            twodee.Rectangle
	Showing           bool
	Frame             int
	observeShowSplash int
}

const (
	OverlayTitleFrame = 0
	OverlayDeathFrame = 1
	OverlayWinFrame   = 2
)

func NewOverlayLayer(app *Application, game *GameLayer) (layer *OverlayLayer, err error) {
	layer = &OverlayLayer{
		Game:   game,
		Events: app.GameEventHandler,
		Bounds: twodee.Rect(0, 0, 1, 1),
	}
	tilem := twodee.TileMetadata{
		Path:       "assets/overlays.fw.png",
		PxPerUnit:  320,
		TileWidth:  320,
		TileHeight: 320,
		FramesWide: 2,
		FramesHigh: 2,
	}
	if layer.TileRenderer, err = twodee.NewTileRenderer(layer.Bounds, app.WinBounds, tilem); err != nil {
		return
	}
	layer.observeShowSplash = layer.Events.AddObserver(ShowSplash, layer.OnShowSplash)
	return
}

func (l *OverlayLayer) OnShowSplash(e twodee.GETyper) {
	if event, ok := e.(*ShowSplashEvent); ok {
		l.Show(event.Frame)
	}
}

func (l *OverlayLayer) Delete() {
	if l.TileRenderer != nil {
		l.TileRenderer.Delete()
	}
	l.Events.RemoveObserver(ShowSplash, l.observeShowSplash)
}

func (l *OverlayLayer) Show(frame int) {
	l.Game.Level.Pause()
	l.Showing = true
	l.Frame = frame
}

func (l *OverlayLayer) Render() {
	if !l.Showing {
		return
	}
	l.TileRenderer.Bind()
	l.TileRenderer.Draw(l.Frame, 0.5, 0.5, 0, false, false)
	l.TileRenderer.Unbind()
}

func (l *OverlayLayer) HandleEvent(evt twodee.Event) bool {
	if !l.Showing {
		return true
	}
	switch event := evt.(type) {
	case *twodee.KeyEvent:
		if event.Type != twodee.Press {
			break
		}
		switch event.Code {
		case twodee.KeyUp:
			fallthrough
		case twodee.KeyDown:
			fallthrough
		case twodee.KeyLeft:
			fallthrough
		case twodee.KeyRight:
			fallthrough
		case twodee.KeyEscape:
			fallthrough
		case twodee.KeySpace:
			fallthrough
		case twodee.KeyEnter:
			l.Game.Level.Unpause()
			l.Showing = false
			if l.Frame == OverlayDeathFrame || l.Frame == OverlayWinFrame {
				l.Events.Enqueue(twodee.NewBasicGameEvent(GameIsClosing))
			}
			return false
		}
	}
	return true
}

func (l *OverlayLayer) Update(elapsed time.Duration) {
}

func (l *OverlayLayer) Reset() (err error) {
	return
}
