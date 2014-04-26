package main

import (
	"time"
)

type Tween interface {
	Update(d time.Duration)
	Done() bool
	Current() float32
}

type tween struct {
	duration time.Duration
	elapsed  time.Duration
}

func (t *tween) Update(d time.Duration) {
	t.elapsed += d
}

func (t *tween) Done() bool {
	return t.elapsed >= t.duration
}

type LinearTween struct {
	tween
	start, end float32
}

func (t *LinearTween) Current() float32 {
	var pct = float32(t.elapsed) / float32(t.duration)
	return pct*(t.end-t.start) + t.start
}

func NewLinearTween(start, end float32, duration time.Duration) Tween {
	return &LinearTween{
		start: start,
		end:   end,
		tween: tween{duration: duration, elapsed: 0},
	}
}
