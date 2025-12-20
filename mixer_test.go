package main

import (
	"testing"
)

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
