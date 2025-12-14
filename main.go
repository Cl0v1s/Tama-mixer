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
		s := Place(body, bodyparts)
		Save("generated", s)
	}
}
