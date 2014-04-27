package main

import (
	"../libs/twodee"
	"time"
)

type OverlayLayer struct {
	Game         *GameLayer
	TileRenderer *twodee.TileRenderer
	Bounds       twodee.Rectangle
	Showing      bool
	Frame        int
}

const (
	OverlayTitleFrame = 0
	OverlayDeathFrame = 1
	OverlayWinFrame   = 2
)

func NewOverlayLayer(app *Application, game *GameLayer) (layer *OverlayLayer, err error) {
	layer = &OverlayLayer{
		Game:   game,
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
	return
}

func (l *OverlayLayer) Delete() {
	if l.TileRenderer != nil {
		l.TileRenderer.Delete()
	}
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
		l.Game.Level.Unpause()
		l.Showing = false
		return false
	}
	return true
}

func (l *OverlayLayer) Update(elapsed time.Duration) {
}

func (l *OverlayLayer) Reset() (err error) {
	return
}
