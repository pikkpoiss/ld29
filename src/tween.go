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
	Start, End float32
}

func (t *LinearTween) Current() float32 {
	var pct = float32(t.elapsed) / float32(t.duration)
	return pct*(t.End-t.Start) + t.Start
}

func NewLinearTween(start, end float32, duration time.Duration) *LinearTween {
	return &LinearTween{
		Start: start,
		End:   end,
		tween: tween{duration: duration, elapsed: 0},
	}
}
