package main

import (
	"runtime"
	"time"

	twodee "../libs/twodee"
	"github.com/go-gl/gl"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

type Application struct {
	layers                *twodee.Layers
	counter               *twodee.Counter
	Context               *twodee.Context
	AudioSystem           *AudioSystem
	WinBounds             twodee.Rectangle
	GameEventHandler      *twodee.GameEventHandler
	gameClosingObserverId int
	InitiateCloseGame     bool
}

func NewApplication() (app *Application, err error) {
	var (
		layers            *twodee.Layers
		context           *twodee.Context
		winbounds         = twodee.Rect(0, 0, 600, 600)
		counter           = twodee.NewCounter()
		gameLayer         *GameLayer
		gameEventHandler  = twodee.NewGameEventHandler(NumGameEventTypes)
		initiateCloseGame = false
	)
	if context, err = twodee.NewContext(); err != nil {
		return
	}
	context.SetFullscreen(false)
	context.SetCursor(false)
	if err = context.CreateWindow(int(winbounds.Max.X), int(winbounds.Max.Y), "LD29"); err != nil {
		return
	}
	layers = twodee.NewLayers()
	app = &Application{
		layers:            layers,
		counter:           counter,
		Context:           context,
		WinBounds:         winbounds,
		GameEventHandler:  gameEventHandler,
		InitiateCloseGame: initiateCloseGame,
	}
	if app.AudioSystem, err = NewAudioSystem(app); err != nil {
		return
	}
	if gameLayer, err = NewGameLayer(app); err != nil {
		return
	}
	if err = gameLayer.LoadLevel("assets/level00/"); err != nil {
		return
	}
	layers.Push(gameLayer)
	app.gameClosingObserverId = app.GameEventHandler.AddObserver(GameIsClosing, app.CloseGame)
	return
}

func (a *Application) Draw() {
	a.counter.Incr()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	a.layers.Render()
}

func (a *Application) Update(elapsed time.Duration) {
	a.layers.Update(elapsed)
}

func (a *Application) Delete() {
	a.GameEventHandler.RemoveObserver(GameIsClosing, a.gameClosingObserverId)
	a.layers.Delete()
	a.AudioSystem.Delete()
	a.Context.Delete()
}

func (a *Application) ProcessEvents() {
	var (
		evt  twodee.Event
		loop = true
	)
	for loop {
		select {
		case evt = <-a.Context.Events.Events:
			a.layers.HandleEvent(evt)
		default:
			// No more events
			loop = false
		}
	}
}

func (a *Application) CloseGame(e twodee.GETyper) {
	a.InitiateCloseGame = true
}

func main() {
	var (
		app *Application
		err error
	)

	if app, err = NewApplication(); err != nil {
		panic(err)
	}
	defer app.Delete()

	var (
		current_time = time.Now()
		updated_to   = current_time
		step         = twodee.Step60Hz
	)
	for !app.Context.Window.ShouldClose() && !app.InitiateCloseGame {
		for !updated_to.After(current_time) {
			app.Update(step)
			updated_to = updated_to.Add(step)
		}
		app.Draw()
		app.Context.Window.SwapBuffers()
		app.Context.Events.Poll()
		app.GameEventHandler.Poll()
		app.ProcessEvents()
		current_time = time.Now()
	}
}
