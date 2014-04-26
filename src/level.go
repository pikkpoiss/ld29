package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	twodee "../libs/twodee"
	"github.com/kurrik/tmxgo"
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

func LoadLevel(path string, eventSystem *twodee.GameEventHandler) (l *Level, err error) {
	var (
		data     []byte
		m        *tmxgo.Map
		layer    tmxgo.Layer
		i, j     int
		grids    = []*twodee.Grid{}
		grid     *twodee.Grid
		maptiles []*tmxgo.Tile
		textiles []twodee.TexturedTile
		maptile  *tmxgo.Tile
		tilemeta twodee.TileMetadata
		batch    *twodee.Batch
		batches  = []*twodee.Batch{}
		player   *Player
	)
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if m, err = tmxgo.ParseMapString(string(data)); err != nil {
		return
	}
	if path, err = GetTexturePath(m, path); err != nil {
		return
	}
	tilemeta = twodee.TileMetadata{
		Path:      path,
		PxPerUnit: 32,
	}
	for i = len(m.Layers) - 1; i >= 0; i-- {
		layer = m.Layers[i]
		if !strings.HasPrefix(layer.Name, "layer") {
			continue
		}
		if maptiles, err = m.TilesFromLayerIndex(int32(i)); err != nil {
			return
		}
		grid = twodee.NewGrid(m.Width, m.Height)
		for j, maptile = range maptiles {
			if maptile != nil {
				grid.SetIndex(int32(j), true)
			}
		}
		grids = append(grids, grid)
		textiles = make([]twodee.TexturedTile, len(maptiles))
		for j, maptile = range maptiles {
			if maptile != nil {
				textiles[j] = maptile
			}
		}
		if batch, err = twodee.LoadBatch(textiles, tilemeta); err != nil {
			return
		}
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

func GetTexturePath(m *tmxgo.Map, path string) (out string, err error) {
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
