package main

import (
	twodee "../libs/twodee"
	"fmt"
	"github.com/kurrik/tmxgo"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Level struct {
	Grids               []*twodee.Grid
	Geometry            []*twodee.Batch
	GridRatios          []float32
	Layers              int
	Player              *Player
	Active              int32
	eventSystem         *twodee.GameEventHandler
	onPlayerMoveEventId int
}

func LoadLevel(path string, names []string, eventSystem *twodee.GameEventHandler) (l *Level, err error) {
	var (
		player  *Player
		grid    *twodee.Grid
		batch   *twodee.Batch
		grids   = []*twodee.Grid{}
		batches = []*twodee.Batch{}
		ratio   float32
		ratios  = []float32{}
	)
	for _, name := range names {
		if grid, batch, ratio, err = loadLayer(path, name); err != nil {
			return
		}
		grids = append(grids, grid)
		batches = append(batches, batch)
		ratios = append(ratios, ratio)
	}
	player = NewPlayer(twodee.NewBaseEntity(1, 1, 1, 1, 0, 0))
	l = &Level{
		Grids:       grids,
		Geometry:    batches,
		Layers:      len(grids),
		GridRatios:  ratios,
		Active:      0,
		Player:      player,
		eventSystem: eventSystem,
	}
	l.onPlayerMoveEventId = eventSystem.AddObserver(PlayerMove, l.OnPlayerMoveEvent)
	return
}

func loadLayer(path, name string) (grid *twodee.Grid, batch *twodee.Batch, ratio float32, err error) {
	var (
		tilemeta twodee.TileMetadata
		maptiles []*tmxgo.Tile
		textiles []twodee.TexturedTile
		maptile  *tmxgo.Tile
		m        *tmxgo.Map
		i        int
		data     []byte
	)
	path = filepath.Join(filepath.Dir(path), name)
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if m, err = tmxgo.ParseMapString(string(data)); err != nil {
		return
	}
	if path, err = getTexturePath(m, path); err != nil {
		return
	}
	tilemeta = twodee.TileMetadata{
		Path:      path,
		PxPerUnit: 32,
	}
	if maptiles, err = m.TilesFromLayerName("collision"); err != nil {
		return
	}
	grid = twodee.NewGrid(m.Width, m.Height)
	for i, maptile = range maptiles {
		if maptile != nil {
			grid.SetIndex(int32(i), true)
		}
	}
	if maptiles, err = m.TilesFromLayerName("tiles"); err != nil {
		return
	}
	textiles = make([]twodee.TexturedTile, len(maptiles))
	for i, maptile = range maptiles {
		if maptile != nil {
			textiles[i] = maptile
		}
	}
	if batch, err = twodee.LoadBatch(textiles, tilemeta); err != nil {
		return
	}
	ratio = float32(grid.Width) * float32(tilemeta.PxPerUnit) / float32(m.TileWidth*m.Width)
	return
}

func getTexturePath(m *tmxgo.Map, path string) (out string, err error) {
	var prefix = filepath.Dir(path)
	for i := 0; i < len(m.Tilesets); i++ {
		if m.Tilesets[i].Image == nil {
			continue
		}
		out = filepath.Join(prefix, m.Tilesets[i].Image.Source)
		return
	}
	err = fmt.Errorf("Could not find suitable tileset")
	return
}

func (l *Level) Delete() {
	for i := 0; i < l.Layers; i++ {
		l.Geometry[i].Delete()
	}
	l.eventSystem.RemoveObserver(PlayerMove, l.onPlayerMoveEventId)
}

func (l *Level) OnPlayerMoveEvent(e twodee.GETyper) {
	if move, ok := e.(*PlayerMoveEvent); ok {
		l.Player.DesiredMove = move.Dir
	}
}

func (l *Level) GridAlignedX(layer int32, p twodee.Point) twodee.Point {
	var (
		ratio = l.GridRatios[layer]
		x     = int32(p.X*ratio + 0.5)
	)
	return twodee.Pt(float32(x)/ratio, p.Y)
}

func (l *Level) GridAlignedY(layer int32, p twodee.Point) twodee.Point {
	var (
		ratio = l.GridRatios[layer]
		y     = int32(p.Y*ratio + 0.5)
	)
	return twodee.Pt(p.X, float32(y)/ratio)
}

func (l *Level) FrontierCollides(layer int32, a, b twodee.Point) bool {
	var (
		ratio = l.GridRatios[layer]
		xmin  = int32(a.X * ratio)
		xmax  = int32(b.X * ratio)
		ymin  = int32(a.Y * ratio)
		ymax  = int32(b.Y * ratio)
	)
	// fmt.Printf("X %v-%v, Y %v-%v\n", xmin, xmax, ymin, ymax)
	for x := xmin; x <= xmax; x++ {
		for y := ymin; y <= ymax; y++ {
			// fmt.Printf("Checking X %v Y %v\n", x, y)
			if l.Grids[layer].Get(x, y) == true {
				return true
			}
		}
	}
	return false
}

func (l *Level) Update(elapsed time.Duration) {
	l.Player.Update(elapsed)
	l.Player.AttemptMove(l)
}
