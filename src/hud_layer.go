package main

import (
	"../libs/twodee"
	"time"
)

type HudLayer struct {
	TileRenderer *twodee.TileRenderer
	Bounds       twodee.Rectangle
	App          *Application
	game         *GameLayer
}

const (
	HudHeight   = 20
	HudWidth    = 20
	HudXPadding = 0
	HudYPadding = 1
	EmptyTile   = 40
	FullTile    = 41
)

func NewHudLayer(app *Application, game *GameLayer) (layer *HudLayer, err error) {
	layer = &HudLayer{
		App:    app,
		Bounds: twodee.Rect(0, 0, HudWidth, HudHeight),
		game:   game,
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
	return
}

func (l *HudLayer) Delete() {
	if l.TileRenderer != nil {
		l.TileRenderer.Delete()
	}
}

func (l *HudLayer) Render() {
	l.TileRenderer.Bind()
	var (
		tiles  = HudHeight - (2 * HudYPadding)
		filled = int(l.game.Level.GetTotalWaterPercent() * float32(tiles))
	)
	for i := 0; i < tiles; i++ {
		var (
			x    = float32(HudWidth-HudXPadding) - 0.5
			y    = float32(i+HudYPadding) + 0.5
			tile = EmptyTile
		)
		if i < filled {
			tile = FullTile
		}
		l.TileRenderer.Draw(tile, x, y, 0, false, false)
	}
	l.TileRenderer.Unbind()
}

func (l *HudLayer) HandleEvent(evt twodee.Event) bool {
	return true
}

func (l *HudLayer) Update(elapsed time.Duration) {
}

func (l *HudLayer) Reset() (err error) {
	return
}
