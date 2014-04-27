package main

import twodee "../libs/twodee"

type AudioSystem struct {
	app                      *Application
	outdoorMusic             *twodee.Music
	exploreMusic             *twodee.Music
	warningMusic             *twodee.Music
	dangerMusic              *twodee.Music
	menuMoveEffect           *twodee.SoundEffect
	menuSelectEffect         *twodee.SoundEffect
	dryWalkEffect            *twodee.SoundEffect
	wetWalkEffect            *twodee.SoundEffect
	outdoorMusicObserverId   int
	exploreMusicObserverId   int
	warningMusicObserverId   int
	dangerMusicObserverId    int
	pauseMusicObserverId     int
	resumeMusicObserverId    int
	menuPauseMusicObserverId int
	menuMoveObserverId       int
	menuSelectObserverId     int
	dryWalkObserverId        int
	wetWalkObserverId        int
	musicToggle              int32
}

func (a *AudioSystem) PlayOutdoorMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
		a.outdoorMusic.Play(-1)
	}
}

func (a *AudioSystem) PlayExploreMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
		a.exploreMusic.Play(-1)
	}
}

func (a *AudioSystem) PlayWarningMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
		a.warningMusic.Play(-1)
	}
}

func (a *AudioSystem) PlayDangerMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
		a.dangerMusic.Play(-1)
	}
}

func (a *AudioSystem) PauseMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
	}
}

func (a *AudioSystem) ResumeMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPaused() {
			twodee.ResumeMusic()
		}
	}
}

func (a *AudioSystem) MenuPauseMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
}

func (a *AudioSystem) PlayMenuMoveEffect(e twodee.GETyper) {
	a.menuMoveEffect.Play(1)
}

func (a *AudioSystem) PlayMenuSelectEffect(e twodee.GETyper) {
	a.menuSelectEffect.Play(1)
}

func (a *AudioSystem) PlayDryWalkEffect(e twodee.GETyper) {
	a.dryWalkEffect.PlayChannel(6, 1)
}

func (a *AudioSystem) PlayWetWalkEffect(e twodee.GETyper) {
	a.wetWalkEffect.PlayChannel(6, 1)
}

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(PlayOutdoorMusic, a.outdoorMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayExploreMusic, a.exploreMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayWarningMusic, a.warningMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayDangerMusic, a.dangerMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PauseMusic, a.pauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ResumeMusic, a.resumeMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(MenuPauseMusic, a.menuPauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(MenuMove, a.menuMoveObserverId)
	a.app.GameEventHandler.RemoveObserver(MenuSelect, a.menuSelectObserverId)
	a.app.GameEventHandler.RemoveObserver(DryWalk, a.dryWalkObserverId)
	a.app.GameEventHandler.RemoveObserver(WetWalk, a.wetWalkObserverId)
	a.outdoorMusic.Delete()
	a.exploreMusic.Delete()
	a.warningMusic.Delete()
	a.dangerMusic.Delete()
	a.menuMoveEffect.Delete()
	a.menuSelectEffect.Delete()
	a.dryWalkEffect.Delete()
	a.wetWalkEffect.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		outdoorMusic     *twodee.Music
		exploreMusic     *twodee.Music
		warningMusic     *twodee.Music
		dangerMusic      *twodee.Music
		menuMoveEffect   *twodee.SoundEffect
		menuSelectEffect *twodee.SoundEffect
		dryWalkEffect    *twodee.SoundEffect
		wetWalkEffect    *twodee.SoundEffect
	)
	if outdoorMusic, err = twodee.NewMusic("assets/music/Outdoor_Theme.ogg"); err != nil {
		return
	}
	if exploreMusic, err = twodee.NewMusic("assets/music/Exploration_Theme.ogg"); err != nil {
		return
	}
	if warningMusic, err = twodee.NewMusic("assets/music/Warning_Theme.ogg"); err != nil {
		return
	}
	if dangerMusic, err = twodee.NewMusic("assets/music/Underwater_Theme.ogg"); err != nil {
		return
	}
	if menuMoveEffect, err = twodee.NewSoundEffect("assets/soundeffects/MenuMove.ogg"); err != nil {
		return
	}
	if menuSelectEffect, err = twodee.NewSoundEffect("assets/soundeffects/MenuSelect.ogg"); err != nil {
		return
	}
	if dryWalkEffect, err = twodee.NewSoundEffect("assets/soundeffects/DryWalk.ogg"); err != nil {
		return
	}
	if wetWalkEffect, err = twodee.NewSoundEffect("assets/soundeffects/WetWalk.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		app:              app,
		outdoorMusic:     outdoorMusic,
		exploreMusic:     exploreMusic,
		warningMusic:     warningMusic,
		dangerMusic:      dangerMusic,
		menuMoveEffect:   menuMoveEffect,
		menuSelectEffect: menuSelectEffect,
		dryWalkEffect:    dryWalkEffect,
		wetWalkEffect:    wetWalkEffect,
		musicToggle:      1,
	}
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayOutdoorMusic, audioSystem.PlayOutdoorMusic)
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayExploreMusic, audioSystem.PlayExploreMusic)
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayWarningMusic, audioSystem.PlayWarningMusic)
	audioSystem.exploreMusicObserverId = app.GameEventHandler.AddObserver(PlayDangerMusic, audioSystem.PlayDangerMusic)
	audioSystem.pauseMusicObserverId = app.GameEventHandler.AddObserver(PauseMusic, audioSystem.PauseMusic)
	audioSystem.resumeMusicObserverId = app.GameEventHandler.AddObserver(ResumeMusic, audioSystem.ResumeMusic)
	audioSystem.menuPauseMusicObserverId = app.GameEventHandler.AddObserver(MenuPauseMusic, audioSystem.MenuPauseMusic)
	audioSystem.menuMoveObserverId = app.GameEventHandler.AddObserver(MenuMove, audioSystem.PlayMenuMoveEffect)
	audioSystem.menuSelectObserverId = app.GameEventHandler.AddObserver(MenuSelect, audioSystem.PlayMenuSelectEffect)
	audioSystem.dryWalkObserverId = app.GameEventHandler.AddObserver(DryWalk, audioSystem.PlayDryWalkEffect)
	audioSystem.wetWalkObserverId = app.GameEventHandler.AddObserver(WetWalk, audioSystem.PlayWetWalkEffect)
	return
}
