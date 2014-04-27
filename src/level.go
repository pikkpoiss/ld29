package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	twodee "../libs/twodee"
	"github.com/kurrik/tmxgo"
)

type Level struct {
	Height                      float32
	Grids                       []*twodee.Grid
	Items                       [][]*Item
	Geometry                    []*twodee.Batch
	GridRatios                  []float32
	Layers                      int32
	Player                      *Player
	Active                      int32
	Transitions                 []*LinearTween
	eventSystem                 *twodee.GameEventHandler
	onPlayerMoveEventId         int
	onPlayerPickedUpItemEventId int
}

func LoadLevel(path string, names []string, eventSystem *twodee.GameEventHandler) (l *Level, err error) {
	var player = NewPlayer(2, 2)
	l = &Level{
		Height:      0,
		Grids:       []*twodee.Grid{},
		Items:       [][]*Item{},
		Geometry:    []*twodee.Batch{},
		GridRatios:  []float32{},
		Layers:      0,
		Active:      0,
		Player:      player,
		eventSystem: eventSystem,
	}
	l.onPlayerMoveEventId = eventSystem.AddObserver(PlayerMove, l.OnPlayerMoveEvent)
	l.onPlayerPickedUpItemEventId = eventSystem.AddObserver(PlayerPickedUpItem, l.OnPlayerPickedUpItemEvent)
	for _, name := range names {
		if err = l.loadLayer(path, name); err != nil {
			return
		}
	}
	return
}

const PxPerUnit float32 = 16.0

func (l *Level) loadLayer(path, name string) (err error) {
	var (
		tilemeta twodee.TileMetadata
		maptiles []*tmxgo.Tile
		textiles []twodee.TexturedTile
		maptile  *tmxgo.Tile
		m        *tmxgo.Map
		i        int
		data     []byte
		height   float32
		grid     *twodee.Grid
		items    []*Item
		batch    *twodee.Batch
		ratio    float32
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
		PxPerUnit: int(PxPerUnit),
	}
	if maptiles, err = m.TilesFromLayerName("entities"); err != nil {
		return
	}
	for i, maptile = range maptiles {
		if maptile != nil {
			items = append(items, NewItem(
				ItemType(maptile.Index),
				"item",
				(maptile.TileBounds.X + maptile.TileBounds.W)/PxPerUnit,
				(maptile.TileBounds.Y + maptile.TileBounds.H)/PxPerUnit,
				maptile.TileBounds.W/PxPerUnit,
				maptile.TileBounds.H/PxPerUnit,
			))
		}
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
	ratio = float32(grid.Width) * PxPerUnit / float32(m.TileWidth*m.Width)
	height = float32(grid.Height) / ratio
	if l.Height < height {
		l.Height = height
	}
	l.Grids = append(l.Grids, grid)
	l.Items = append(l.Items, items)
	l.Geometry = append(l.Geometry, batch)
	l.Layers += 1
	l.Transitions = append(l.Transitions, nil)
	l.GridRatios = append(l.GridRatios, ratio)
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
	var i int32
	for i = 0; i < l.Layers; i++ {
		l.Geometry[i].Delete()
	}
	l.eventSystem.RemoveObserver(PlayerMove, l.onPlayerMoveEventId)
	l.eventSystem.RemoveObserver(PlayerPickedUpItem, l.onPlayerPickedUpItemEventId)
}

func (l *Level) OnPlayerMoveEvent(e twodee.GETyper) {
	if move, ok := e.(*PlayerMoveEvent); ok {
		l.Player.DesiredMove = move.Dir
	}
}

func (l *Level) OnPlayerPickedUpItemEvent(e twodee.GETyper) {
	if pickup, ok := e.(*PlayerPickedUpItemEvent); ok {
		l.Player.AddToInventory(pickup.Item)
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

// Given points a,b defining the leading edge of a moving entity; determine
// if there is a collision with something on the grid.
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
	playerBounds := l.Player.Bounds()
	for _, item := range l.Items[layer] {
		if playerBounds.Overlaps(item.Bounds()) {
			l.eventSystem.Enqueue(NewPlayerPickedUpItemEvent(item))
			break
		}
	}
	return false
}

func (l *Level) GetLayerY(index int32) float32 {
	var tween = l.Transitions[index]
	if tween != nil {
		return tween.Current()
	}
	switch {
	case index > l.Active:
		return -1
	case index < l.Active:
		return l.Height
	case index == l.Active:
		return 0
	}
	return 0
}

const TopSlideSpeed = time.Duration(320) * time.Millisecond
const BotSlideSpeed = time.Duration(320) * time.Millisecond

func (l *Level) LayerAdvance() {
	if l.Active >= l.Layers-1 {
		return
	}
	l.Transitions[l.Active] = NewLinearTween(0, l.Height, TopSlideSpeed)
	l.Active++
	l.Transitions[l.Active] = NewLinearTween(-1, 0, BotSlideSpeed)
}

func (l *Level) LayerRewind() {
	if l.Active <= 0 {
		return
	}
	l.Transitions[l.Active-1] = NewLinearTween(l.Height, 0, TopSlideSpeed)
	l.Transitions[l.Active-1].SetCallback(func() {
		l.Active--
	})
	l.Transitions[l.Active] = NewLinearTween(0, -1, BotSlideSpeed)
}

func (l *Level) Update(elapsed time.Duration) {
	var i int32
	for i = 0; i < l.Layers; i++ {
		if l.Transitions[i] != nil {
			if l.Transitions[i].Done() {
				l.Transitions[i] = nil
			} else {
				l.Transitions[i].Update(elapsed)
			}
		}
	}
	l.Player.Update(elapsed)
	l.Player.AttemptMove(l)
}
