package main

import (
	"math"
	"testing"
)

const simplePath = "M 130 10 C 120 20, 180 20, 170 10"
const intermediatePath = "m 10 10 c 1 2 3 4 5 6 z m -10 -10"
const complexPath = " m 0 0 c 2.384706 -3.8247189 9.090522 -2.8303014 13.419508 -2.9698914 4.328986 -0.13959 8.777591 -0.1227708 12.099557 1.8699314 3.321966 1.992702 5.665494 5.552875 6.599759 8.799678 0.934265 3.246803 0.907 6.146902 -0.769972 9.239661 -1.401888 2.585435 -4.005582 4.719819 -7.320097 6.355104 -0.650388 0.320882 -1.328144 0.622547 -2.029561 0.904631 -4.275988 1.719648 -11.763687 4.446312 -15.729423 1.649939 -0.37043 -0.261202 -0.689892 -0.548786 -0.965528 -0.859054 -2.675262 -3.011389 -1.222003 -8.159713 -2.169357 -12.065472 -0.706428 -2.912467 -2.980282 -6.045882 -3.64267 -8.952141 -0.317429 -1.392734 -0.264783 -2.733301 0.507784 -3.972386 z"

func TestParseDSimple(t *testing.T) {
	test := simplePath
	commands := ParseD(test)
	if len(commands) != 2 {
		t.Error("len(commands) must be 2")
	}
	if commands[0].Type != "M" || commands[1].Type != "C" {
		t.Error("Wrong command types")
	}
}

func TestParseDIntermediate(t *testing.T) {
	test := intermediatePath
	commands := ParseD(test)
	if len(commands) != 4 {
		t.Error("len(commands) must be 4")
	}
	if commands[0].Type != "M" || commands[1].Type != "C" || commands[2].Type != "Z" || commands[3].Type != "M" {
		t.Error("Wrong command types")
	}
	if commands[0].Args[0] != 10 || commands[0].Args[1] != 10 || commands[3].Args[0] != 0 || commands[3].Args[1] != 0 {
		t.Error("Wrong M Args")
	}
	if commands[1].Args[0] != 11 || commands[1].Args[1] != 12 || commands[1].Args[2] != 13 || commands[1].Args[3] != 14 || commands[1].Args[4] != 15 || commands[1].Args[5] != 16 {
		t.Error("Wrong C Args")
	}

}

func TestParseDComplex(t *testing.T) {
	test := complexPath
	commands := ParseD(test)
	if len(commands) != 3 {
		t.Error("len(commands) must be 3")
	}
	if commands[0].Type != "M" || commands[1].Type != "C" || commands[2].Type != "Z" {
		t.Error("Wrong command types")
	}
}

func TestGetBeziersFromCommandsSimple(t *testing.T) {
	beziers := GetBeziersFromCommands(ParseD(simplePath))
	if len(beziers) != 1 {
		t.Error("len(beziers) must be 1")
	}
}

func TestGetBeziersFromCommandsZ(t *testing.T) {
	const pathWithZ = "M 10 10 C 120 20, 180 20, 170 10 Z"
	got := GetBeziersFromCommands(ParseD(pathWithZ))
	wantStart := Point{X: 170, Y: 10}
	wantEnd := Point{X: 10, Y: 10}
	if len(got) != 2 {
		t.Error("len(got) must be 2")
	}
	if got[1].P0.X != wantStart.X && got[1].P0.Y != wantStart.Y {
		t.Error("got[1].P0 is different from wantStart")
	}
	if got[1].P3.X != wantEnd.X && got[1].P3.Y != wantEnd.Y {
		t.Error("got[1].P3 is different from wantEnd")
	}
}

func TestGetBeziersFromCommandsComplex(t *testing.T) {
	beziers := GetBeziersFromCommands(ParseD(complexPath))
	if len(beziers) != 12 {
		t.Errorf("len(beziers) must be 12, it is %d", len(beziers))
	}
}

func TestGetPointFromBezier(t *testing.T) {
	bezier := Bezier{
		P0: Point{X: 0, Y: 0},
		P1: Point{X: 0, Y: 10},
		P2: Point{X: 10, Y: 10},
		P3: Point{X: 10, Y: 0},
	}

	tests := []struct {
		t        float64
		expected Point
	}{
		{0, Point{X: 0, Y: 0}},
		{0.5, Point{X: 5, Y: 7.5}},
		{1, Point{X: 10, Y: 0}},
	}

	for _, tc := range tests {
		p := GetPointFromBezier(bezier, tc.t)
		if math.Abs(p.X-tc.expected.X) > 1e-9 || math.Abs(p.Y-tc.expected.Y) > 1e-9 {
			t.Errorf("t=%.2f → got (%.3f, %.3f), expected (%.3f, %.3f)",
				tc.t, p.X, p.Y, tc.expected.X, tc.expected.Y)
		}
	}
}

func TestFindClosestPointInPathSimple(t *testing.T) {
	paths := []Path{{D: "M 130 10 C 120 20, 180 20, 170 10"}}
	got, _ := findClosestPointInPaths(
		paths,
		Point{X: 130, Y: 10},
		1,
	)
	want := 0.0
	if got != want {
		t.Errorf("Got %f expected %f", got, want)
	}
	got, _ = findClosestPointInPaths(
		paths,
		Point{X: 170, Y: 10},
		1,
	)
	want = 1.0
	if got != want {
		t.Errorf("Got %f expected %f", got, want)
	}
}

func TestFindClosestPointInPathComplex(t *testing.T) {
	paths := []Path{{D: "m 0 0 c 2.384706 -3.8247189 9.090522 -2.8303014 13.419508 -2.9698914 4.328986 -0.13959 8.777591 -0.1227708 12.099557 1.8699314"}}

	got, _ := findClosestPointInPaths(
		paths,
		Point{X: 0, Y: 0},
		1,
	)
	want := 0.0
	if got != want {
		t.Errorf("Got %f expected %f", got, want)
	}
	got, _ = findClosestPointInPaths(
		paths,
		Point{X: 0 + 13.419508 + 12.099557, Y: 0 - 2.9698914 + 1.8699314},
		1,
	)
	want = 1.0
	if got != want {
		t.Errorf("Got %f expected %f", got, want)
	}
}

func TestFindClosestPointInPathError(t *testing.T) {
	paths := []Path{{D: "M 130 10 C 120 20, 180 20, 170 10"}}
	got, _ := findClosestPointInPaths(
		paths,
		Point{X: 0, Y: 10},
		1,
	)
	want := -1.0
	if got != want {
		t.Errorf("Got %f expected %f", got, want)
	}
}
