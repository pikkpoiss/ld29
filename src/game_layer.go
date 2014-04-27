package main

import (
	twodee "../libs/twodee"
	"os"
	"strings"
	"time"
)

type GameLayer struct {
	Level         *Level
	BatchRenderer *twodee.BatchRenderer
	TileRenderer  *twodee.TileRenderer
	Bounds        twodee.Rectangle
	App           *Application
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{
		App:    app,
		Bounds: twodee.Rect(0, 0, 20, 20),
	}
	if layer.BatchRenderer, err = twodee.NewBatchRenderer(layer.Bounds, app.WinBounds); err != nil {
		return
	}
	tilem := twodee.TileMetadata{
		Path:       "assets/entities.fw.png",
		PxPerUnit:  int(PxPerUnit),
		TileWidth:  32,
		TileHeight: 32,
	}
	if layer.TileRenderer, err = twodee.NewTileRenderer(layer.Bounds, app.WinBounds, tilem); err != nil {
		return
	}
	layer.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayExploreMusic))
	return
}

func (l *GameLayer) LoadLevel(path string) (err error) {
	var (
		dir        *os.File
		candidates []string
		names      = []string{}
	)
	if dir, err = os.Open(path); err != nil {
		return
	}
	if candidates, err = dir.Readdirnames(0); err != nil {
		return
	}
	for _, name := range candidates {
		if strings.HasPrefix(name, "layer") {
			names = append(names, name)
		}
	}
	l.Level, err = LoadLevel(path, names, l.App.GameEventHandler)
	return
}

func (l *GameLayer) Delete() {
	if l.Level != nil {
		l.Level.Delete()
	}
}

func (l *GameLayer) Render() {
	l.BatchRenderer.Bind()
	var i int32
	var y float32
	for i = l.Level.Layers - 1; i >= 0; i-- {
		switch {
		case i == l.Level.Active:
			fallthrough
		case i == l.Level.Active+1:
			fallthrough
		case i == l.Level.Active-1:
			y = l.Level.GetLayerY(i)
			l.BatchRenderer.Draw(l.Level.Geometry[i], 0, y, 0)
		}
		if i == l.Level.Active {
			l.BatchRenderer.Unbind()
			l.TileRenderer.Bind()

			for _, item := range l.Level.Items[i] {
				pt := item.Pos()
				l.TileRenderer.Draw(item.Frame(), pt.X, pt.Y+y, 0, false, false)
			}

			pt := l.Level.Player.Pos()
			l.TileRenderer.Draw(l.Level.Player.Frame(), pt.X, pt.Y+y, 0, l.Level.Player.FlippedX(), false)

			l.TileRenderer.Unbind()
			l.BatchRenderer.Bind()
		}
	}
	l.BatchRenderer.Unbind()
	return
}

func (l *GameLayer) Update(elapsed time.Duration) {
	l.Level.Update(elapsed)
}

func (l *GameLayer) Reset() (err error) {
	return
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	var released bool
	switch event := evt.(type) {
	case *twodee.KeyEvent:
		released = event.Type == twodee.Release
		switch event.Code {
		case twodee.KeyUp:
			if released {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(None))
			} else {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(North))
			}
		case twodee.KeyRight:
			if released {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(None))
			} else {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(East))
			}
		case twodee.KeyDown:
			if released {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(None))
			} else {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(South))
			}
		case twodee.KeyLeft:
			if released {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(None))
			} else {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(West))
			}
		case twodee.KeyJ:
			if released {
				l.Level.LayerRewind()
			}
		case twodee.KeyK:
			if released {
				l.Level.LayerAdvance()
			}
		}
	}
	return true
}
