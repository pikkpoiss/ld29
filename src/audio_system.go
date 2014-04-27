package main

import twodee "../libs/twodee"

type AudioSystem struct {
	app                    *Application
	outdoorMusic           *twodee.Music
	exploreMusic           *twodee.Music
	warningMusic           *twodee.Music
	dangerMusic            *twodee.Music
	outdoorMusicObserverId int
	exploreMusicObserverId int
	warningMusicObserverId int
	dangerMusicObserverId  int
	pauseMusicObserverId   int
	resumeMusicObserverId  int
}

func (a *AudioSystem) PlayOutdoorMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
	a.outdoorMusic.Play(-1)
}

func (a *AudioSystem) PlayExploreMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
	a.exploreMusic.Play(-1)
}

func (a *AudioSystem) PlayWarningMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
	a.warningMusic.Play(-1)
}

func (a *AudioSystem) PlayDangerMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
	a.dangerMusic.Play(-1)
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
	a.app.GameEventHandler.RemoveObserver(PlayOutdoorMusic, a.outdoorMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayExploreMusic, a.exploreMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayWarningMusic, a.warningMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayDangerMusic, a.dangerMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PauseMusic, a.pauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ResumeMusic, a.resumeMusicObserverId)
	a.exploreMusic.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		outdoorMusic *twodee.Music
		exploreMusic *twodee.Music
		warningMusic *twodee.Music
		dangerMusic  *twodee.Music
	)
	if outdoorMusic, err = twodee.NewMusic("assets/music/Outdoor_Theme_1.ogg"); err != nil {
		return
	}
	if exploreMusic, err = twodee.NewMusic("assets/music/Exploration_Theme_1.ogg"); err != nil {
		return
	}
	if warningMusic, err = twodee.NewMusic("assets/music/Warning_Theme_1.ogg"); err != nil {
		return
	}
	if dangerMusic, err = twodee.NewMusic("assets/music/Danger_Theme_1.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		app:          app,
		outdoorMusic: outdoorMusic,
		exploreMusic: exploreMusic,
		warningMusic: warningMusic,
		dangerMusic:  dangerMusic,
	}
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayOutdoorMusic, audioSystem.PlayOutdoorMusic)
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayExploreMusic, audioSystem.PlayExploreMusic)
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayWarningMusic, audioSystem.PlayWarningMusic)
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayDangerMusic, audioSystem.PlayDangerMusic)
	audioSystem.pauseMusicObserverId = app.GameEventHandler.AddObserver(PauseMusic, audioSystem.PauseMusic)
	audioSystem.resumeMusicObserverId = app.GameEventHandler.AddObserver(ResumeMusic, audioSystem.ResumeMusic)
	return
}
