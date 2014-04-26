package main

import (
	twodee "../libs/twodee"
	"fmt"
	"github.com/kurrik/tmxgo"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Level struct {
	Grids []*twodee.Grid
}

func LoadLevel(path string) (l *Level, err error) {
	var (
		data  []byte
		m     *tmxgo.Map
		layer tmxgo.Layer
		i, j  int
		grids = []*twodee.Grid{}
		grid  *twodee.Grid
		tiles []*tmxgo.Tile
		tile  *tmxgo.Tile
	)
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if m, err = tmxgo.ParseMapString(string(data)); err != nil {
		return
	}
	for i, layer = range m.Layers {
		if !strings.HasPrefix(layer.Name, "layer") {
			continue
		}
		if tiles, err = m.TilesFromLayerIndex(int32(i)); err != nil {
			return
		}
		grid = twodee.NewGrid(m.Width, m.Height)
		for j, tile = range tiles {
			if tile != nil {
				grid.SetIndex(int32(j), true)
			}
		}
		grids = append(grids, grid)
	}
	l = &Level{
		Grids: grids,
	}
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
