package main

import (
	"encoding/xml"
	"os"
)

func main() {
	file, err := os.Open("tama.svg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var svg SVG
	decoder := xml.NewDecoder(file)

	if err := decoder.Decode(&svg); err != nil {
		panic(err)
	}

	bodies, bodyparts := Sort(svg)

	for i := 0; i < len(bodyparts); i++ {
		x, y := FindLowerPadding(bodyparts[i], bodyparts[i].Label)
		bodyparts[i] = Unpad(bodyparts[i], x, y)
		Save("bodyparts", []Group{bodyparts[i]}, -20, 20)
	}

	for i := 0; i < len(bodies); i++ {
		x, y := FindLowerPadding(bodies[i], "body")
		bodies[i] = Unpad(bodies[i], x, y)
		Save("bodies", []Group{bodies[i]}, -20, 20)
		bodies[i] = Place(bodies[i], bodyparts)
		Save("generated", []Group{bodies[i]}, -20, 20)
	}
}
