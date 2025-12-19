package main

import (
	"math"
	"testing"
)

const simplePath = "M 130 10 C 120 20, 180 20, 170 10"
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

func TestParseDComplex(t *testing.T) {
	test := complexPath
	commands := ParseD(test)
	if len(commands) != 3 {
		t.Error("len(commands) must be 3")
	}
	if commands[0].Type != "m" || commands[1].Type != "c" || commands[2].Type != "z" {
		t.Error("Wrong command types")
	}
}

func TestGetBeziersFromCommandsSimple(t *testing.T) {
	beziers := GetBeziersFromCommands(ParseD(simplePath))
	if len(beziers) != 1 {
		t.Error("len(beziers) must be 1")
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

func TestGetRotationFromBezier(t *testing.T) {
	bezier := Bezier{
		P0: Point{X: 0, Y: 0},
		P1: Point{X: 1, Y: 0},
		P2: Point{X: 2, Y: 0},
		P3: Point{X: 3, Y: 0},
	}
	got := GetRotationFromBezier(bezier, 0.5)
	want := 0.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("got %v, want %v", got, want)
	}

	bezier = Bezier{
		P0: Point{X: 0, Y: 0},
		P1: Point{X: 0, Y: 1},
		P2: Point{X: 0, Y: 2},
		P3: Point{X: 0, Y: 3},
	}
	got = GetRotationFromBezier(bezier, 0.5)
	want = 90.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("got %v, want %v", got, want)
	}
}
