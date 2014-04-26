package main

import (
	twodee "../libs/twodee"
	"github.com/kurrik/tmxgo"
	"strings"
)

type Level struct {
	Grids []*twodee.Grid
}

func LoadLevel(path string) (l *Level, err error) {
	var (
		data  []byte
		m     *tmxgo.Map
		layer *tmxgo.Layer
		i, j  int32
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
		if tiles, err = m.TilesFromLayerIndex(i); err != nil {
			return
		}
		grid = twodee.NewGrid(m.Width, m.Height)
		for i, t = range tiles {
			if t != nil {
				grid.SetIndex(i, true)
			}
		}
		grids = append(grids, grid)
	}
	l = &Level{
		Grids: grids,
	}
	return
}
