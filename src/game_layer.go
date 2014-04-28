package main

import (
	twodee "../libs/twodee"
	"os"
	"sort"
	"strings"
	"time"
)

func NewRain() *twodee.AnimatingEntity {
	return twodee.NewAnimatingEntity(
		0, 0,
		32.0/PxPerUnit, 32.0/PxPerUnit,
		0,
		twodee.Step10Hz,
		[]int{
			56, 57, 58, 59, 60, 61, 62,
			64, 65, 66, 67, 68, 69, 70,
		},
	)
}

func NewWater() *twodee.AnimatingEntity {
	return twodee.NewAnimatingEntity(
		0, 0,
		32.0/PxPerUnit, 32.0/PxPerUnit,
		0,
		twodee.Step10Hz,
		[]int{
			72, 73, 74, 75,
		},
	)
}

func NewPump() *twodee.AnimatingEntity {
	return twodee.NewAnimatingEntity(
		0, 0,
		32.0/PxPerUnit, 32.0/PxPerUnit,
		0,
		twodee.Step10Hz,
		[]int{
			44, 45, 46, 45,
		},
	)
}

const (
	WaterEntityOffset = 80
)

type GameLayer struct {
	Level                     *Level
	BatchRenderer             *twodee.BatchRenderer
	TileRenderer              *twodee.TileRenderer
	Bounds                    twodee.Rectangle
	App                       *Application
	menuResumeMusicObserverId int
	Rain                      *twodee.AnimatingEntity
	Water                     *twodee.AnimatingEntity
	Pump                      *twodee.AnimatingEntity
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{
		App:    app,
		Bounds: twodee.Rect(0, 0, 20, 20),
		Rain:   NewRain(),
		Water:  NewWater(),
		Pump:   NewPump(),
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
	layer.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayOutdoorMusic))
	layer.menuResumeMusicObserverId = layer.App.GameEventHandler.AddObserver(MenuResumeMusic, layer.MenuResumeMusic)
	//layer.App.GameEventHandler.Enqueue(NewShowSplashEvent(OverlayDeathFrame))
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
	sort.Sort(sort.StringSlice(candidates))
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
	if l.TileRenderer != nil {
		l.TileRenderer.Delete()
	}
	if l.BatchRenderer != nil {
		l.BatchRenderer.Delete()
	}
	l.App.GameEventHandler.RemoveObserver(MenuResumeMusic, l.menuResumeMusicObserverId)
}

func (l *GameLayer) Render() {
	l.BatchRenderer.Bind()
	var i int32
	var y float32
	for i = l.Level.Layers - 1; i >= 0; i-- {
		var waterStatus = l.Level.GetLayerWaterStatus(i)
		switch {
		case i == l.Level.Active:
			fallthrough
		case i == l.Level.Active+1:
			fallthrough
		case i == l.Level.Active-1:
			y = l.Level.GetLayerY(i)
			switch waterStatus {
			case Dry:
				l.Level.Geometry[i].SetTextureOffsetPx(0, 0)
			case Wet:
				l.Level.Geometry[i].SetTextureOffsetPx(0, -16)
			case Flooded:
				l.Level.Geometry[i].SetTextureOffsetPx(0, 0)
			}
			l.BatchRenderer.Draw(l.Level.Geometry[i], 0, y, 0)
		}
		if i == l.Level.Active {
			l.BatchRenderer.Unbind()
			l.TileRenderer.Bind()

			for _, item := range l.Level.Items[i] {
				pt := item.Pos()
				frame := item.Frame()
				if item.Id == ItemPump && l.Level.Player.IsPumping {
					frame = l.Pump.Frame()
				}
				if waterStatus == Wet && i != 0 {
					frame += WaterEntityOffset
				}
				l.TileRenderer.Draw(frame, pt.X, pt.Y+y, 0, false, false)
			}

			pt := l.Level.Player.Pos()
			frame := l.Level.Player.Frame()
			if waterStatus == Wet && i != 0 {
				frame += WaterEntityOffset
			}
			l.TileRenderer.Draw(frame, pt.X, pt.Y+y, 0, l.Level.Player.FlippedX(), false)
			if i == 0 {
				l.DrawOverlay(l.Rain)
			} else if waterStatus == Flooded {
				l.DrawOverlay(l.Water)
			}
			l.TileRenderer.Unbind()
			l.BatchRenderer.Bind()
		}
	}
	l.BatchRenderer.Unbind()
	return
}

func (l *GameLayer) DrawOverlay(entity *twodee.AnimatingEntity) {
	var j = 0
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			l.TileRenderer.Draw(entity.OffsetFrame(j), float32(2*x)+1, float32(2*y)+1, 0, false, false)
			j += 3
		}
	}
}

func (l *GameLayer) Update(elapsed time.Duration) {
	l.Level.Update(elapsed)
	l.Rain.Update(elapsed)
	l.Water.Update(elapsed)
	l.Pump.Update(elapsed)
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
				l.App.GameEventHandler.Enqueue(NewInversePlayerMoveEvent(North))
			} else {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(North))
			}
		case twodee.KeyRight:
			if released {
				l.App.GameEventHandler.Enqueue(NewInversePlayerMoveEvent(East))
			} else {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(East))
			}
		case twodee.KeyDown:
			if released {
				l.App.GameEventHandler.Enqueue(NewInversePlayerMoveEvent(South))
			} else {
				l.App.GameEventHandler.Enqueue(NewPlayerMoveEvent(South))
			}
		case twodee.KeyLeft:
			if released {
				l.App.GameEventHandler.Enqueue(NewInversePlayerMoveEvent(West))
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

func (l *GameLayer) MenuResumeMusic(e twodee.GETyper) {
	var layerWaterStatus = l.Level.GetLayerWaterStatus(l.Level.Active)
	if l.Level.Active == 0 {
		l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayOutdoorMusic))
	} else if layerWaterStatus == 0 {
		l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayExploreMusic))
	} else if layerWaterStatus == 1 {
		l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayWarningMusic))
	} else if layerWaterStatus == 2 {
		l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayDangerMusic))
	}
}
