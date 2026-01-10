package main

import (
	"encoding/xml"
	"fmt"
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

	for _, bodypart := range bodyparts {
		Save("out/bodyparts", bodypart.Svg)
	}

	for _, body := range bodies {
		Save("out/bodies", body.Svg)
	}

	fmt.Println("Mixing")
	bodies = Mix(bodies, bodyparts)
	fmt.Println("Mixing done")
	fmt.Println("Assembling")
	for i := 0; i < len(bodies); i++ {
		// fmt.Printf("%d / %d\n", i, len(bodies))
		bodies[i] = BodyAssemble(bodies[i])
		bodies[i] = BodyReframe(bodies[i], 32)
	}
	fmt.Println("Assembling done")

	for _, body := range bodies {
		Save("out/generated", body.Svg)
	}

	grouppedFrames := ParseFrames("out/generated")

	for _, formGroup := range grouppedFrames {
		for _, expressionGroup := range formGroup {
			SaveFrames("out/generated", "out/sorted", expressionGroup)
		}
	}

}
