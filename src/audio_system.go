package main

import twodee "../libs/twodee"

type AudioSystem struct {
	app                    *Application
	exploreMusic           *twodee.Music
	exploreMusicObserverId int
	pauseMusicObserverId   int
	resumeMusicObserverId  int
}

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
	a.app.GameEventHandler.RemoveObserver(PlayExploreMusic, a.exploreMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PauseMusic, a.pauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ResumeMusic, a.resumeMusicObserverId)
	a.exploreMusic.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
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
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayExploreMusic, audioSystem.PlayExploreMusic)
	audioSystem.pauseMusicObserverId = app.GameEventHandler.AddObserver(PauseMusic, audioSystem.PauseMusic)
	audioSystem.resumeMusicObserverId = app.GameEventHandler.AddObserver(ResumeMusic, audioSystem.ResumeMusic)
	return
}
