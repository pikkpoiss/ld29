package main

import twodee "../libs/twodee"

type AudioSystem struct {
	app *Application

	// TODO: Add in *twodee.Music and *twodee.SoundEffect pointers here
	// for adutio elements being used
	exploreMusic           *twodee.Music
	exploreMusicObserverId int
}

// TODO: Add in functions to play music files and sound effect files
func (a *AudioSystem) PlayExploreMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
	a.exploreMusic.Play(-1)
}

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
	a.exploreMusic.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	// TODO: Set up vars for mMsic and sound effect assests
	// Create new Music and Sound Effect objects
	var (
		exploreMusic *twodee.Music
	)
	if exploreMusic, err = twodee.NewMusic("assets/music/Journey_Theme_1.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		app:          app,
		exploreMusic: exploreMusic,
	}
	// TODO: Set up observers for each music and sound effect related
	// function
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayExploreMusic, audioSystem.PlayExploreMusic)
	return
}
