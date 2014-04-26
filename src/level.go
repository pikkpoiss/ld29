package main

import (
	twodee "../libs/twodee"
	"github.com/kurrik/tmxgo"
)

type Level struct {
}

func LoadLevel(path string) (l *Level, err error) {
	var (
		data []byte
		m    *tmxgo.Map
	)
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if m, err = tmxgo.ParseMapString(string(data)); err != nil {
		return
	}
}
