package main

import twodee "../libs/twodee"

type AudioSystem struct {
	app *Application

	// TODO: Add in *twodee.Music and *twodee.SoundEffect pointers here
	// for adutio elements being used
}

// TODO: Add in functions to play music files and sound effect files

func (a *AudioSystem) PauseMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
}

func (a *AudioSystem) ResumeMusic(e twodee.GETyper) {
	if twodee.MusicIsPaused() {
		twodee.ResumeMusic()
	}
}

func (a *AudioSystem) Delete() {
	// TODO: Call RemoveObserver and Delete functions on Game Event
	// Handlers and Music/Sound Effect objects respectively
}

func NewAudioSystem(app *Application) (ausioSystem *AudioSystem, err error) {
	// TODO: Set up vars for mMsic and sound effect assests
	// Create new Music and Sound Effect objects

	audioSystem = &AudioSystem{
		app: app,
	}

	// TODO: Set up observers for each music and sound effect related
	// function
}
