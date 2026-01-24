package main

import (
	"testing"
)

func TestCommandTransformSimple(t *testing.T) {
	test := "M 0 0"
	tr := Transformation{Translation: Point{X: 10, Y: 0}}
	commands := ParseD(test)
	commands[0].Transform(tr)
	if commands[0].Args[0] != 10 || commands[0].Args[1] != 0 {
		t.Errorf("Bad transformed coordinates")
	}
}
