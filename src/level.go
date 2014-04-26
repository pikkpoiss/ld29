package main

import (
	twodee "../libs/twodee"
	"fmt"
	"github.com/kurrik/tmxgo"
	"io/ioutil"
	"path/filepath"
)

type Level struct {
	Grids               []*twodee.Grid
	Geometry            []*twodee.Batch
	Layers              int
	Player              *Player
	eventSystem         *twodee.GameEventHandler
	onPlayerMoveEventId int
}

func (l *Level) OnPlayerMoveEvent(e twodee.GETyper) {
	fmt.Println("Got player move event.")
}

func LoadLevel(path string, names []string, eventSystem *twodee.GameEventHandler) (l *Level, err error) {
	var (
		player  *Player
		grid    *twodee.Grid
		batch   *twodee.Batch
		grids   = []*twodee.Grid{}
		batches = []*twodee.Batch{}
	)
	for _, name := range names {
		if grid, batch, err = loadLayer(path, name); err != nil {
			return
		}
		grids = append(grids, grid)
		batches = append(batches, batch)
	}
	player = NewPlayer(twodee.NewBaseEntity(1, 1, 32, 32, 0, 0))
	l = &Level{
		Grids:       grids,
		Geometry:    batches,
		Layers:      len(grids),
		Player:      player,
		eventSystem: eventSystem,
	}
	l.onPlayerMoveEventId = eventSystem.AddObserver(PlayerMove, l.OnPlayerMoveEvent)
	return
}

func loadLayer(path, name string) (grid *twodee.Grid, batch *twodee.Batch, err error) {
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
