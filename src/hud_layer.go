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
	HudHeight         = 20
	HudWidth          = 20
	HudWaterXPadding  = 0
	HudWaterYPadding  = 1
	HudHealthXPadding = 2
	HudHealthYPadding = 0
	HudItemXPadding   = 1
	HudItemYPadding   = 1.5
	EmptyWaterTile    = 40
	FullWaterTile     = 41
	EmptyHealthTile   = 42
	FullHealthTile    = 43
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

func (l *HudLayer) RenderWater() {
	var (
		tiles  = HudHeight - (2 * HudWaterYPadding)
		filled = int(l.game.Level.GetTotalWaterPercent() * float32(tiles))
	)
	for i := 0; i < tiles; i++ {
		var (
			x    = float32(HudWidth-HudWaterXPadding) - 0.5
			y    = float32(i+HudWaterYPadding) + 0.5
			tile = EmptyWaterTile
		)
		if i < filled {
			tile = FullWaterTile
		}
		l.TileRenderer.Draw(tile, x, y, 0, false, false)
	}
}

func (l *HudLayer) RenderHealth() {
	var (
		health = l.game.Level.Player.HealthPercent()
		tiles  = HudWidth - (2 * HudHealthXPadding)
		filled = int(health * float32(tiles))
	)
	if health >= 1.0 {
		return
	}
	for i := 0; i < tiles; i++ {
		var (
			x    = float32(i+HudHealthXPadding) + 0.5
			y    = float32(HudHeight-HudHealthYPadding) - 0.5
			tile = EmptyHealthTile
		)
		if i < filled {
			tile = FullHealthTile
		}
		l.TileRenderer.Draw(tile, x, y, 0, false, false)
	}
}

func (l *HudLayer) RenderItems() {
	for i, item := range l.game.Level.Player.Inventory {
		var (
			x = float32(i+HudItemXPadding) + 0.5
			y = float32(HudItemYPadding) - 0.5
		)
		l.TileRenderer.Draw(item.Frame(), x, y, 0, false, false)
	}
}

func (l *HudLayer) Render() {
	l.TileRenderer.Bind()
	l.RenderWater()
	l.RenderHealth()
	l.RenderItems()
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
