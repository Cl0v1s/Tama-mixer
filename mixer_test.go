package main

import (
	"testing"
)

func TestBezierToDSimple(t *testing.T) {
	// M 130 10 C 120 20, 180 20, 170 10
	bezier := Bezier{P0: Point{X: 130, Y: 10}, P1: Point{X: 120, Y: 20}, P2: Point{X: 180, Y: 20}, P3: Point{X: 170, Y: 10}}
	got := bezierToD([]Bezier{bezier, bezier})
	want := " M 130 10 C 120 20 180 20 170 10  M 130 10 C 120 20 180 20 170 10 "
	if got != want {
		t.Errorf("Got %s expected %s", got, want)
	}
}
