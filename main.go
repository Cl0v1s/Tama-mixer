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

	for _, bodypart := range bodyparts {
		Save("bodyparts", bodypart.Svg)
	}

	for _, body := range bodies {
		Save("bodies", body.Svg)
	}

	bodies = Mix(bodies, bodyparts)
	for i := 0; i < len(bodies); i++ {
		bodies[i] = BodyAssemble(bodies[i])
	}

	for _, body := range bodies {
		Save("generated", body.Svg)
	}

}
