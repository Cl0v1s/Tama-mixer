package main

import (
	"encoding/xml"
	"os"
)

func main() {
	file, err := os.Open("svg/parts.svg")
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
	bodiesGroups := GroupBodies(bodies)
	bodypartsGroups := GroupBodyParts(bodyparts)
	for _, group := range bodypartsGroups {
		SaveBodyPartsToJSON("out/bodyparts/", group)
	}
	for _, group := range bodiesGroups {
		SaveBodiesToJSON("out/bodies/", group)
	}
}
