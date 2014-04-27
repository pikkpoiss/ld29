package main

import (
	"image/color"
	"time"

	twodee "../libs/twodee"
)

const (
	ProgramCode int32 = iota
)

const (
	ExitCode int32 = iota
	MusicCode
)

type MenuLayer struct {
	visible        bool
	menu           *twodee.Menu
	text           *twodee.TextRenderer
	regularFont    *twodee.FontFace
	cache          map[int]*twodee.TextCache
	highlightCache *twodee.TextCache
	activeCache    *twodee.TextCache
	bounds         twodee.Rectangle
	App            *Application
}

func NewMenuLayer(window twodee.Rectangle, App *Application) (layer *MenuLayer, err error) {
	var (
		menu          *twodee.Menu
		regularFont   *twodee.FontFace
		highlightFont *twodee.FontFace
		activeFont    *twodee.FontFace
		background    = color.Transparent
		font          = "assets/fonts/slkscr.ttf"
	)
	if regularFont, err = twodee.NewFontFace(font, 32, color.RGBA{200, 200, 200, 255}, background); err != nil {
		return
	}
	if highlightFont, err = twodee.NewFontFace(font, 32, color.RGBA{255, 240, 120, 255}, background); err != nil {
		return
	}
	if activeFont, err = twodee.NewFontFace(font, 32, color.RGBA{200, 200, 255, 255}, background); err != nil {
		return
	}
	menu, err = twodee.NewMenu([]twodee.MenuItem{
		//twodee.NewParentMenuItem("Music", []twodee.MenuItem{
		//	twodee.NewBackMenuItem(".."),
		//	twodee.NewBoundValueMenuItem("On", 1, &App.AudioSystem.musicToggle),
		//	twodee.NewBoundValueMenuItem("Off", 0, &App.AudioSystem.musicToggle),
		//}),
		twodee.NewKeyValueMenuItem("Music On/Off", ProgramCode, MusicCode),
		twodee.NewKeyValueMenuItem("Exit", ProgramCode, ExitCode),
	})
	if err != nil {
		return
	}
	layer = &MenuLayer{
		menu:           menu,
		regularFont:    regularFont,
		cache:          map[int]*twodee.TextCache{},
		activeCache:    twodee.NewTextCache(activeFont),
		highlightCache: twodee.NewTextCache(highlightFont),
		bounds:         window,
		App:            App,
		visible:        false,
	}
	err = layer.Reset()
	return
}

func (l *MenuLayer) Reset() (err error) {
	if l.text != nil {
		l.text.Delete()
	}
	if l.text, err = twodee.NewTextRenderer(l.bounds); err != nil {
		return
	}
	l.activeCache.Clear()
	l.highlightCache.Clear()
	for _, v := range l.cache {
		v.Clear()
	}
	return
}

func (l *MenuLayer) Delete() {
	l.text.Delete()
	l.activeCache.Delete()
	l.highlightCache.Delete()
	for _, v := range l.cache {
		v.Delete()
	}
}

func (l *MenuLayer) Render() {
	if !l.visible {
		return
	}
	var (
		textCache *twodee.TextCache
		texture   *twodee.Texture
		ok        bool
		y         = l.bounds.Max.Y
	)
	l.text.Bind()
	for i, item := range l.menu.Items() {
		if item.Highlighted() {
			l.highlightCache.SetText(item.Label())
			texture = l.highlightCache.Texture
		} else if item.Active() {
			l.activeCache.SetText(item.Label())
			texture = l.activeCache.Texture
		} else {
			if textCache, ok = l.cache[i]; !ok {
				textCache = twodee.NewTextCache(l.regularFont)
				l.cache[i] = textCache
			}
			textCache.SetText(item.Label())
			texture = textCache.Texture
		}
		if texture != nil {
			y = y - float32(texture.Height)
			l.text.Draw(texture, 0, y)
		}
	}
	l.text.Unbind()
}

func (l *MenuLayer) Update(elapsed time.Duration) {
}

func (l *MenuLayer) HandleEvent(evt twodee.Event) bool {
	if !l.visible {
		switch event := evt.(type) {
		case *twodee.KeyEvent:
			if event.Type != twodee.Press {
				break
			}
			if event.Code == twodee.KeyEscape {
				l.menu.Reset()
				l.visible = true
				l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuSelect))
			}
		}
		return true
	}
	switch event := evt.(type) {
	case *twodee.MouseButtonEvent:
		if event.Type != twodee.Press {
			break
		}
		if data := l.menu.Select(); data != nil {
			l.handleMenuItem(data)
		}
		l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuSelect))
	case *twodee.MouseMoveEvent:
		var (
			y         = l.bounds.Max.Y
			my        = y - event.Y
			texture   *twodee.Texture
			textCache *twodee.TextCache
			ok        bool
		)
		for i, item := range l.menu.Items() {
			if item.Highlighted() {
				texture = l.highlightCache.Texture
			} else if item.Active() {
				texture = l.activeCache.Texture
			} else {
				if textCache, ok = l.cache[i]; ok {
					texture = textCache.Texture
				}
			}
			if texture != nil {
				y = y - float32(texture.Height)
				if my >= y {
					if !item.Highlighted() {
						l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuMove))
						l.menu.HighlightItem(item)
					}
					break
				}
			}
		}
	case *twodee.KeyEvent:
		if event.Type != twodee.Press {
			break
		}
		switch event.Code {
		case twodee.KeyEscape:
			l.visible = false
			l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuSelect))
			return false
		case twodee.KeyUp:
			l.menu.Prev()
			l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuMove))
			return false
		case twodee.KeyDown:
			l.menu.Next()
			l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuMove))
			return false
		case twodee.KeyEnter:
			if data := l.menu.Select(); data != nil {
				l.handleMenuItem(data)
			}
			l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuSelect))
			return false
		}
	}
	return true
}

func (l *MenuLayer) handleMenuItem(data *twodee.MenuItemData) {
	switch data.Key {
	case ProgramCode:
		switch data.Value {
		case MusicCode:
			if twodee.MusicIsPaused() {
				l.App.AudioSystem.musicToggle = 1
				l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuResumeMusic))
			} else {
				l.App.AudioSystem.musicToggle = 0
				l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(MenuPauseMusic))
			}
		case ExitCode:
			l.App.InitiateCloseGame = true
		}
	}
}
