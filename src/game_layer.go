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
	playerLayer   int
	topLayer      int
	topTrans      *LinearTween
	botLayer      int
	botTrans      *LinearTween
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{
		App:      app,
		Bounds:   twodee.Rect(0, 0, 10, 10),
		topLayer: 0,
		botLayer: 1,
	}
	if layer.BatchRenderer, err = twodee.NewBatchRenderer(layer.Bounds, app.WinBounds); err != nil {
		return
	}
	tilem := twodee.TileMetadata{
		Path:       "assets/entities.fw.png",
		PxPerUnit:  32,
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
	var (
		y float32
	)
	l.BatchRenderer.Bind()
	for i := l.Level.Layers - 1; i >= 0; i-- {
		y = 0.0
		if i == l.topLayer {
			if l.topTrans != nil {
				y = l.topTrans.Current()
			} else if l.topLayer != l.playerLayer {
				y = l.Bounds.Max.Y
			}
		} else if i == l.botLayer {
			if l.botTrans != nil {
				y = l.botTrans.Current()
			} else if l.botLayer != l.playerLayer {
				y = -1.0
			}
		}
		if i == l.topLayer || i == l.botLayer {
			l.BatchRenderer.Draw(l.Level.Geometry[i], 0, y, 0)
		}
		if i == l.playerLayer {
			l.BatchRenderer.Unbind()
			l.TileRenderer.Bind()
			pt := l.Level.Player.Pos()
			l.TileRenderer.Draw(l.Level.Player.Frame(), pt.X, pt.Y, 0, false, false)
			l.TileRenderer.Unbind()
			l.BatchRenderer.Bind()
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
	l.Level.Update(elapsed)
}

const TopSlideSpeed = time.Duration(500) * time.Millisecond
const BotSlideSpeed = time.Duration(200) * time.Millisecond

func (l *GameLayer) LayerAdvance() {
	if l.playerLayer >= l.Level.Layers-1 || l.Level == nil {
		return
	}
	l.topLayer = l.playerLayer
	l.botLayer = l.topLayer + 1
	l.playerLayer = l.botLayer
	l.topTrans = NewLinearTween(0, l.Bounds.Max.Y, TopSlideSpeed)
	l.botTrans = NewLinearTween(-1, 0, BotSlideSpeed)
}

func (l *GameLayer) LayerRewind() {
	if l.playerLayer <= 0 || l.Level == nil {
		return
	}
	l.topLayer = l.playerLayer - 1
	l.botLayer = l.topLayer + 1
	l.playerLayer = l.topLayer
	l.topTrans = NewLinearTween(l.Bounds.Max.Y, 0, TopSlideSpeed)
	l.botTrans = NewLinearTween(0, -1, BotSlideSpeed)
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
		case twodee.KeyEscape:
			l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(GameIsClosing))
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
				l.LayerRewind()
			}
		case twodee.KeyK:
			if released {
				l.LayerAdvance()
			}
		case twodee.KeyM:
			if twodee.MusicIsPaused() {
				l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(ResumeMusic))
			} else {
				l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PauseMusic))
			}
		}
	}
	return true
}
