package main

import (
	"fmt"
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

func TestBezierToDComplexe(t *testing.T) {
	d := "m -0.3454999999999977 -2.1071930000000005 c 0 0 -4.49345 -1.318558 -6.578528 -0.507048 -0.843157 0.328156 -2.666533 1.310354 -2.024162 2.028188 1.234839 1.379903 8.836053 0.638468 8.836053 0.638468 0 0 -7.798679 1.086497 -8.330013 1.643245 -0.707129 0.740952 1.382069 1.72642 2.313329 1.94368 2.183934 0.509505 6.506237 -1.774665 6.506237 -1.774665 z"
	commands := ParseD(d)
	beziers := GetBeziersFromCommands(commands)
	got := bezierToD(beziers)
	fmt.Println(got)
}
