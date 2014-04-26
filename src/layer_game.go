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
		if i == l.topLayer || i == l.botLayer || i == l.playerLayer {
			l.BatchRenderer.Draw(l.Level.Geometry[i], 0, y, 0)
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
	switch event := evt.(type) {
	case *twodee.KeyEvent:
		if event.Type == twodee.Release {
			break
		}
		switch event.Code {
		case twodee.KeyEscape:
			l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(GameIsClosing))
		case twodee.KeyUp:
			l.LayerRewind()
			l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(North))
		case twodee.KeyRight:
			l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(East))
		case twodee.KeyDown:
			l.LayerAdvance()
			l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(South))
		case twodee.KeyLeft:
			l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(West))
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
