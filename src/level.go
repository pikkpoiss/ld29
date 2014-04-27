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
	Height                       float32
	Grids                        []*twodee.Grid
	Items                        [][]*Item
	Geometry                     []*twodee.Batch
	GridRatios                   []float32
	Layers                       int32
	Player                       *Player
	Active                       int32
	Transitions                  []*LinearTween
	eventSystem                  *twodee.GameEventHandler
	onPlayerMoveEventId          int
	onPlayerTouchedItemEventId   int
	onPlayerDestroyedItemEventId int
	onPlayerUsedItemEventId      int
	WaterAccumulation            time.Duration
	Paused                       bool
}

func LoadLevel(path string, names []string, eventSystem *twodee.GameEventHandler) (l *Level, err error) {
	var player = NewPlayer(10, 5)
	l = &Level{
		Height:            0,
		Grids:             []*twodee.Grid{},
		Items:             [][]*Item{},
		Geometry:          []*twodee.Batch{},
		GridRatios:        []float32{},
		Layers:            0,
		Active:            0,
		Player:            player,
		eventSystem:       eventSystem,
		WaterAccumulation: 0,
		Paused:            false,
	}
	l.onPlayerMoveEventId = eventSystem.AddObserver(PlayerMove, l.OnPlayerMoveEvent)
	l.onPlayerTouchedItemEventId = eventSystem.AddObserver(PlayerTouchedItem, l.OnPlayerTouchedItemEvent)
	l.onPlayerUsedItemEventId = eventSystem.AddObserver(PlayerUsedItem, l.OnPlayerUsedItemEvent)
	l.onPlayerDestroyedItemEventId = eventSystem.AddObserver(PlayerDestroyedItem, l.OnPlayerDestroyedItemEvent)
	for _, name := range names {
		if err = l.loadLayer(path, name); err != nil {
			return
		}
	}
	return
}

const LevelWaterThreshold time.Duration = time.Duration(30) * time.Second

type LayerWaterStatus int

const (
	Dry LayerWaterStatus = iota
	Wet
	Flooded
)

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
			itemId := ItemId(maptile.Index)
			items = append(items, NewItem(
				itemId,
				ItemIdToType[itemId],
				"item",
				(maptile.TileBounds.X+maptile.TileBounds.W)/PxPerUnit,
				(maptile.TileBounds.Y+maptile.TileBounds.H)/PxPerUnit,
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
	//batch.SetTextureOffsetPx(0, 16)
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
	l.eventSystem.RemoveObserver(PlayerTouchedItem, l.onPlayerTouchedItemEventId)
	l.eventSystem.RemoveObserver(PlayerDestroyedItem, l.onPlayerDestroyedItemEventId)
	l.eventSystem.RemoveObserver(PlayerUsedItem, l.onPlayerUsedItemEventId)
}

func (l *Level) OnPlayerMoveEvent(e twodee.GETyper) {
	if move, ok := e.(*PlayerMoveEvent); ok {
		l.Player.UpdateDesiredMove(move.Dir, move.Inverse)
		//		l.Player.DesiredMove = move.Dir
	}
}

func (l *Level) OnPlayerTouchedItemEvent(e twodee.GETyper) {
	if !l.Player.CanGetItem {
		return
	}
	if touched, ok := e.(*PlayerTouchedItemEvent); ok {
		l.Player.CanGetItem = false
		switch touched.Item.Type {
		case LayerThresholdItem:
			l.Player.MoveTo(touched.Item.Pos())
			l.Player.CanMove = false
			switch touched.Item.Id {
			case ItemUp:
				l.LayerRewind()
			case ItemDown:
				l.LayerAdvance()
			}
		case UseableItem:
			l.Player.MoveTo(touched.Item.Pos())
			l.Player.CanMove = false
			l.eventSystem.Enqueue(NewPlayerUsedItemEvent(touched.Item))
		case InventoryItem:
			l.RemoveItem(touched.Item)
			l.Player.AddToInventory(touched.Item)
		case DestructableItem:
			if l.Player.CanDestroy(touched.Item) {
				l.eventSystem.Enqueue(NewPlayerDestroyedItemEvent(touched.Item))
			}
		}
	}
}

func (l *Level) OnPlayerUsedItemEvent(e twodee.GETyper) {
	if used, ok := e.(*PlayerUsedItemEvent); ok {
		switch used.Item.Id {
		case ItemPump:
			l.Player.IsPumping = true
		}
	}
}

func (l *Level) OnPlayerDestroyedItemEvent(e twodee.GETyper) {
	if destroyed, ok := e.(*PlayerDestroyedItemEvent); ok {
		l.RemoveItem(destroyed.Item)
	}
}

// Removes the item from the current layer's Items slice.
func (l *Level) RemoveItem(item *Item) {
	layerItems := l.Items[l.Active]
	index := -1
	for i, levelItem := range layerItems {
		if levelItem == item {
			index = i
			break
		}
	}
	if index != -1 {
		copy(layerItems[index:], layerItems[index+1:])
		layerItems[len(layerItems)-1] = nil
		layerItems = layerItems[:len(layerItems)-1]
		// Be sure to update the slice on the the level.
		l.Items[l.Active] = layerItems
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
	touchedItem := false
	for _, item := range l.Items[layer] {
		if playerBounds.Overlaps(item.Bounds()) {
			l.eventSystem.Enqueue(NewPlayerTouchedItemEvent(item))
			touchedItem = true
			break
		}
	}
	if !touchedItem {
		// Reset the player's pumping state.
		l.Player.IsPumping = false
		// Prevent the player from triggering another item
		// pickup until they've moved off of all items
		l.Player.CanGetItem = true
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
	var newWaterLevel = l.GetLayerWaterStatus(l.Active + 1)
	var previousWaterLevel = l.GetLayerWaterStatus(l.Active)
	if l.Active == 0 {
		if newWaterLevel == 0 {
			l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayExploreMusic))
		} else if newWaterLevel == 1 {
			l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayWarningMusic))
		} else if newWaterLevel == 2 {
			l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayDangerMusic))
		}
	} else if newWaterLevel == 1 && previousWaterLevel == 0 {
		l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayWarningMusic))
	} else if newWaterLevel == 2 && previousWaterLevel != 2 {
		l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayDangerMusic))
	}
	l.Transitions[l.Active] = NewLinearTween(0, l.Height, TopSlideSpeed)
	l.Player.SetState(ClimbDown | Down)
	l.Player.SetCallback(func() {
		l.Player.CanMove = true
		l.Player.SetState(Standing | Down)
	})
	l.Active++
	l.Transitions[l.Active] = NewLinearTween(-1, 0, BotSlideSpeed)
}

func (l *Level) LayerRewind() {
	if l.Active <= 0 {
		return
	}
	var newWaterLevel = l.GetLayerWaterStatus(l.Active - 1)
	var previousWaterLevel = l.GetLayerWaterStatus(l.Active)
	if l.Active == 1 {
		l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayOutdoorMusic))
	} else if newWaterLevel == 0 && previousWaterLevel != 0 {
		l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayExploreMusic))
	} else if newWaterLevel == 1 && previousWaterLevel != 0 {
		l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayWarningMusic))
	}
	l.Transitions[l.Active-1] = NewLinearTween(l.Height, 0, TopSlideSpeed)
	l.Transitions[l.Active-1].SetCallback(func() {
		l.Active--
		l.Player.SetState(ClimbUp | Down)
		l.Player.SetCallback(func() {
			l.Player.CanMove = true
			l.Player.SetState(Standing | Down)
			if l.Active == 0 && l.Player.HasFinalItem {
				l.eventSystem.Enqueue(NewShowSplashEvent(OverlayWinFrame))
			}
		})
	})
	l.Transitions[l.Active] = NewLinearTween(0, -1, BotSlideSpeed)
}

func (l *Level) Update(elapsed time.Duration) {
	if l.Paused {
		return
	}
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
	var currentWaterStatus = l.GetLayerWaterStatus(l.Active)
	l.WaterAccumulation += elapsed
	if l.Player.IsPumping {
		l.WaterAccumulation -= 2 * elapsed
		if l.WaterAccumulation < 0 {
			l.WaterAccumulation = 0
		}
	}
	var newWaterStatus = l.GetLayerWaterStatus(l.Active)
	if l.Active != 0 && (newWaterStatus > currentWaterStatus) {
		if newWaterStatus == 1 {
			l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayWarningMusic))
		} else {
			l.eventSystem.Enqueue(twodee.NewBasicGameEvent(PlayDangerMusic))
		}
	}
	if l.Active != 0 && currentWaterStatus == Flooded {
		l.Player.Damage(PlayerWaterDamage)
	}
	if l.Player.HealthPercent() == 0 {
		l.eventSystem.Enqueue(NewShowSplashEvent(OverlayDeathFrame))
	}
}

func (l *Level) GetTotalWaterPercent() float32 {
	return float32(l.WaterAccumulation) / float32(LevelWaterThreshold)
}

func (l *Level) GetLayerWaterStatus(layer int32) LayerWaterStatus {
	var percentFlooded = l.GetTotalWaterPercent()
	if percentFlooded >= 1 {
		return Flooded
	}
	var layerLevelBottom = 1 - float32(layer)/float32(l.Layers)
	var layerLevelTop = layerLevelBottom + (1.00 / float32(l.Layers))
	if percentFlooded >= layerLevelTop {
		return Flooded
	} else if percentFlooded >= layerLevelBottom {
		return Wet
	}
	return Dry
}

func (l *Level) Pause() {
	l.Paused = true
}

func (l *Level) Unpause() {
	l.Paused = false
}
