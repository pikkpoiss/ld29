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
	layers      *twodee.Layers
	counter     *twodee.Counter
	Context     *twodee.Context
	AudioSystem *AudioSystem
}

func NewApplication() (app *Application, err error) {
	var (
		layers      *twodee.Layers
		context     *twodee.Context
		winbounds   = twodee.Rect(0, 0, 600, 600)
		counter     = twodee.NewCounter()
		audioSystem *AudioSystem
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
		layers:  layers,
		counter: counter,
		Context: context,
	}
	if audioSystem, err = NewAudioSystem(app); err != nil {
		return
	}
	app.AudioSystem = audioSystem
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
	for !app.Context.Window.ShouldClose() {
		for !updated_to.After(current_time) {
			app.Update(step)
			updated_to = updated_to.Add(step)
		}
		app.Draw()
		app.Context.Window.SwapBuffers()
		app.Context.Events.Poll()
		app.ProcessEvents()
		current_time = time.Now()
	}
}
